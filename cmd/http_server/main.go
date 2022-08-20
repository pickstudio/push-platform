package main

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pickstudio/push-platform/constants"

	xrayecs "github.com/aws/aws-xray-sdk-go/awsplugins/ecs"
	"github.com/aws/aws-xray-sdk-go/xray"
	pp "github.com/pickstudio/push-platform"
	handlerhttp "github.com/pickstudio/push-platform/internal/handler/http"
	"github.com/pickstudio/push-platform/pkg/recov"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Netflix/go-env"
	"github.com/rs/zerolog/log"

	oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"
	edgechi "github.com/pickstudio/push-platform/edge/chi"
	edgesqs "github.com/pickstudio/push-platform/edge/sqs"
	adaptermessage "github.com/pickstudio/push-platform/internal/adapter/message"
	"github.com/pickstudio/push-platform/internal/config"
	servicemessage "github.com/pickstudio/push-platform/internal/service/message"
)

var (
	cfg        config.Config
	httpServer *http.Server
)

func init() {
	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Interface("config", cfg).Msg("http_server start")
}

func main() {

	// fatal detection
	defer func() {
		if t := recover(); t != nil {
			err := recov.RecoverFn(t)
			if err != nil {
				log.Error().Err(err).Str("type", "fatal").Msg("anomaly terminated")
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sqs, err := edgesqs.New(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("sqs is not settings up")
	}

	messageAdapter, err := adaptermessage.New(
		ctx, sqs,
		cfg.AWSSQSQueue.Name, cfg.AWSSQSQueue.Timeout,
		cfg.AWSSQSDeadLetterQueue.Name, cfg.AWSSQSDeadLetterQueue.Timeout,
	)
	if err != nil {
		log.Panic().Err(err).Msg("failed to inject package that adapter message.")
	}
	messageService := servicemessage.New(
		messageAdapter,
	)

	handler := handlerhttp.New(messageService)

	r := edgechi.New()

	fsStatic, err := fs.Sub(pp.StaticSwaggerUI, "static/swagger-ui")
	if err != nil {
		log.Panic().Err(err).Msg("serve swagger static files")
	}
	r.Mount("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.FS(fsStatic))))

	pushStatic, err := fs.Sub(pp.OAPISpecYAML, "static/push")
	if err != nil {
		log.Panic().Err(err).Msg("serve push static files")
	}
	r.Mount("/push/", http.StripPrefix("/swagger/", http.FileServer(http.FS(pushStatic))))

	fsSpec, err := fs.Sub(pp.OAPISpecYAML, "api/oapi")
	if err != nil {
		log.Panic().Err(err).Msg("serve api specifications")
	}
	r.Mount("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(fsSpec))))

	oapiv1.HandlerWithOptions(handler, oapiv1.ChiServerOptions{
		BaseURL:          "/api/v1",
		BaseRouter:       r,
		Middlewares:      []oapiv1.MiddlewareFunc{},
		ErrorHandlerFunc: edgechi.OAPIErrorHandler,
	})

	if cfg.Monitoring {
		r.Handle("/metrics", promhttp.Handler())

		xrayecs.Init()
		xray.Handler(xray.NewFixedSegmentNamer(constants.ValueProject), r)
	}

	httpServer = &http.Server{
		Handler: r,
		Addr:    cfg.LocalhostHTTP.DSN,
	}

	go func() {
		if err := http.ListenAndServe(cfg.LocalhostHTTP.DSN, r); err != nil {
			log.Err(err).Msgf("listen: %s", cfg.LocalhostHTTP.DSN)
			panic(err)
		}
	}()
	shutdown(ctx, httpServer)
}

func shutdown(ctx context.Context, srv *http.Server) {
	var stop = make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)
	<-stop

	log.Info().Msg("......... Stopping 'http' server")

	if err := srv.Shutdown(ctx); err != nil {
		log.Err(err).Msgf("http server shutdown")
	}

	log.Info().Msg("......... It maybe shutdown gracefully, GoodBye")
	close(stop)
}
