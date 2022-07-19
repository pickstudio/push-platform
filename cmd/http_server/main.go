package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Netflix/go-env"
	"github.com/rs/zerolog/log"

	oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"
	edgechi "github.com/pickstudio/push-platform/edge/chi"
	"github.com/pickstudio/push-platform/internal/config"
	"github.com/pickstudio/push-platform/internal/handler"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	handler := handler.New()
	r := edgechi.New()
	oapiv1.HandlerWithOptions(handler, oapiv1.ChiServerOptions{
		BaseURL:          "/api/v1",
		BaseRouter:       r,
		Middlewares:      []oapiv1.MiddlewareFunc{},
		ErrorHandlerFunc: edgechi.OAPIErrorHandler,
	})
	httpServer = &http.Server{
		Handler: r,
		Addr:    cfg.LocalhostHttp.DSN,
	}

	go func() {
		if err := http.ListenAndServe(cfg.LocalhostHttp.DSN, r); err != nil {
			log.Err(err).Msgf("listen: %s", cfg.LocalhostHttp.DSN)
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
