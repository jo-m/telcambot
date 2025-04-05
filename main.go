package main

import (
	"image"

	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog/log"
	"jo-m.ch/go/telcambot/pkg/libcamera"
	"jo-m.ch/go/telcambot/pkg/logging"
	"jo-m.ch/go/telcambot/pkg/tel"
)

type config struct {
	logging.LogConfig
	tel.TelConfig
	libcamera.Config
}

func main() {
	c := config{}
	arg.MustParse(&c)
	logging.MustInit(c.LogConfig)

	log.Info().Msg("starting")

	bot, err := tel.NewBot(c.TelConfig, func() (image.Image, error) {
		return libcamera.Snap(c.Config)
	})
	if err != nil {
		log.Panic().Err(err).Send()
	}

	log.Info().Msg("initialized")

	bot.RunForever()
}
