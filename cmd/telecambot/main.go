package main

import (
	"github.com/alexflint/go-arg"
	"github.com/jo-m/telecambot/pkg/logging"
	"github.com/jo-m/telecambot/pkg/tel"
	"github.com/rs/zerolog/log"
)

type config struct {
	logging.LogConfig
	tel.TelConfig

	DBFile string `arg:"env:DB_FILE,--db-file" default:"db.sqlite3" help:"sqlite3 database file" placeholder:"FILE"`
}

func main() {
	c := config{}
	arg.MustParse(&c)
	logging.MustInit(c.LogConfig)

	log.Info().Msg("starting")

	bot, err := tel.NewBot(c.TelConfig)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	log.Info().Msg("initialized")

	bot.RunForever()
}
