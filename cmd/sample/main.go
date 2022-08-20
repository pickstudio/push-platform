package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pickstudio/push-platform/internal/model"
	"github.com/pickstudio/push-platform/pkg/arrays"
	"github.com/pickstudio/push-platform/pkg/recov"

	"github.com/Netflix/go-env"
	"github.com/rs/zerolog/log"

	edgesqs "github.com/pickstudio/push-platform/edge/sqs"
	adaptermessage "github.com/pickstudio/push-platform/internal/adapter/message"
	"github.com/pickstudio/push-platform/internal/config"
	servicemessage "github.com/pickstudio/push-platform/internal/service/message"
)

var (
	cfg config.Config
)

func init() {
	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Interface("config", cfg).Msg("http_server start")
}

/*
	beforeProcessor
	processor,
	afterProcessor,
	fallback,
*/
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sqs, err := edgesqs.New(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("sqs is not settings up")
	}

	messageAdapter, err := adaptermessage.New(
		ctx,
		sqs,
		cfg.AWSSQSQueue.Name, cfg.AWSSQSQueue.Timeout,
		cfg.AWSSQSDeadLetterQueue.Name, cfg.AWSSQSDeadLetterQueue.Timeout,
	)
	if err != nil {
		panic(err.Error())
	}

	messageService := servicemessage.New(
		messageAdapter,
	)

	count, failedMessage, err := messageService.PushMessage(
		ctx,
		[]*model.Message{
			{
				ID:   uuid.NewString(),
				From: "local",

				Service: model.MessageServiceBuddyStock,
				UserID:  uuid.NewString(),

				Device:    model.MessageDeviceAndroid,
				PushToken: uuid.NewString(),

				ViewType: model.MessageViewTypePlain,
				PlainView: &model.PlainView{
					Title:        "와우 안녕",
					Content:      "친구들",
					ThumbnailURL: "https://naver.com",
					SchemeURL:    "buddystock://buddy/12",
					Alarm:        "wow",
					CreatedAt:    time.Now(),
				},
			},
		},
	)
	if err != nil {
		panic(err.Error())
	}
	log.Info().Int("count", count).Interface("failedMessage", failedMessage).Err(err).Send()

	msgList, err := messageService.ReceiveMessage(ctx)
	if err != nil {
		panic(err.Error())
	}
	arrays.ForEach(msgList, func(v *model.Message, _ int) {
		switch v.Service {
		case model.MessageServiceBuddyStock:
			if v.Device == model.MessageDeviceAndroid {
				log.Panic().Msg("buddystock-android")
			} else if v.Device == model.MessageDeviceIOS {
				log.Panic().Msg("buddystock-ios")
			}
			log.Panic().Msg("buddystock-panic")
		case model.MessageServicePickMe:
			if v.Device == model.MessageDeviceAndroid {
				log.Panic().Msg("pickme-android")
			} else if v.Device == model.MessageDeviceIOS {
				log.Panic().Msg("pickme-ios")
			}
			log.Panic().Msg("pickme-panic")
		case model.MessageServiceDijkstra:
			if v.Device == model.MessageDeviceAndroid {
				log.Panic().Msg("dijkstra-android")
			} else if v.Device == model.MessageDeviceIOS {
				log.Panic().Msg("dijkstra-ios")
			}
			log.Panic().Msg("dijkstra-panic")
		}
		log.Panic().Msg("nothings~")
	})

	log.Info().Interface("success_list", msgList).Send()
}
