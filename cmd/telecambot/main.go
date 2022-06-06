package main

import (
	"image"

	"github.com/alexflint/go-arg"
	"github.com/jo-m/telecambot/pkg/logging"
	"github.com/jo-m/telecambot/pkg/pic"
	"github.com/jo-m/telecambot/pkg/tel"
	"github.com/rs/zerolog/log"
)

type config struct {
	logging.LogConfig
	tel.TelConfig
	pic.CamConfig
}

func main() {
	c := config{}
	arg.MustParse(&c)
	logging.MustInit(c.LogConfig)

	log.Info().Msg("starting")

	bot, err := tel.NewBot(c.TelConfig, func() (image.Image, error) {
		return pic.Snap(c.CamConfig)
	})
	if err != nil {
		log.Panic().Err(err).Send()
	}

	log.Info().Msg("initialized")

	bot.RunForever()
}
