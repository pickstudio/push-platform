package main

import (
	"github.com/Netflix/go-env"
	"github.com/pickstudio/push-platform/internal/config"
	"github.com/rs/zerolog/log"
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

func main() {

}
