package tel

import (
	"github.com/rs/zerolog/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelConfig struct {
	BotAPIToken string `arg:"required,env:BOT_API_TOKEN,--bot-api-token" help:"get it from https://t.me/Botfather" placeholder:"TOKEN"`
	BotDebug    bool   `arg:"--bot-debug,env:BOT_DEBUG" default:"false" help:"run telegram bot in debug mode"`
}

type Bot struct {
	config TelConfig
	api    *tgbotapi.BotAPI
}

func NewBot(config TelConfig) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.BotAPIToken)
	if err != nil {
		return nil, err
	}

	api.Debug = config.BotDebug

	return &Bot{
		config: config,
		api:    api,
	}, nil
}

func (b *Bot) RunForever() {
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	for update := range b.api.GetUpdatesChan(update) {
		log.Debug().Interface("update", update).Msg("new update")
	}
}
