package logging

import (
	"log/syslog"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogConfig struct {
	LogPretty bool   `arg:"--log-pretty,env:LOG_PRETTY" default:"false" help:"log pretty"`
	LogLevel  string `arg:"--log-level,env:LOG_LEVEL" default:"info" help:"log level" placeholder:"LEVEL"`
	LogSyslog bool   `arg:"--log-syslog,env:LOG_SYSLOG" default:"false" help:"log to syslog, disables pretty logging"`
}

func MustInit(config LogConfig) {
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	if config.LogSyslog {
		syslogger, err := syslog.New(syslog.LOG_INFO|syslog.LOG_USER, "")
		if err != nil {
			log.Panic().Err(err).Send()
		}
		log.Logger = log.Output(zerolog.SyslogLevelWriter(syslogger))
	} else if config.LogPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"
	zerolog.TimestampFunc = func() time.Time { return time.Now() }

	log.Logger = log.Logger.Level(level).With().Timestamp().Caller().Logger()
}
