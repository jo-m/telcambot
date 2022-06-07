package tel

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"time"

	"github.com/rs/zerolog/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelConfig struct {
	BotAPIToken string  `arg:"required,env:BOT_API_TOKEN,--bot-api-token" help:"get it from https://t.me/Botfather" placeholder:"TOK"`
	BotDebug    bool    `arg:"--bot-debug,env:BOT_DEBUG" default:"false" help:"run telegram bot in debug mode"`
	BotUsers    []int64 `arg:"--bot-users,env:BOT_USERS" help:"whitelist of Telegram user ids allowed to use the bot" placeholder:"IDS,"`
}

type Bot struct {
	config TelConfig
	api    *tgbotapi.BotAPI

	snapFn SnapFn
}

type SnapFn func() (image.Image, error)

func NewBot(config TelConfig, snapFn SnapFn) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.BotAPIToken)
	if err != nil {
		return nil, err
	}

	api.Debug = config.BotDebug

	return &Bot{
		config: config,
		api:    api,

		snapFn: snapFn,
	}, nil
}

func (b *Bot) userAllowed(from tgbotapi.User) bool {
	for _, adminID := range b.config.BotUsers {
		if from.ID == adminID {
			return true
		}
	}

	return false
}

func (b *Bot) replyTo(msg *tgbotapi.Message, text string) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
	newMsg.ReplyToMessageID = msg.MessageID
	_, err := b.api.Send(newMsg)
	if err != nil {
		log.Err(err).Send()
	}
}

func (b *Bot) sendPic(msg *tgbotapi.Message, ts time.Time, im image.Image) error {
	buf := bytes.Buffer{}
	err := jpeg.Encode(&buf, im, nil)
	if err != nil {
		return err
	}

	file := tgbotapi.FileBytes{
		Name:  fmt.Sprintf("%s.jpeg", ts.Format(time.RFC3339)),
		Bytes: buf.Bytes(),
	}

	sendPic := tgbotapi.NewPhoto(msg.Chat.ID, file)
	sendPic.Caption = ts.Format(time.RFC1123)
	_, err = b.api.Send(sendPic)
	return err
}

func (b *Bot) handlePicCommand(msg *tgbotapi.Message) {
	now := time.Now()
	im, err := b.snapFn()
	if err != nil {
		log.Err(err).Send()
		b.replyTo(msg, fmt.Sprintf("Could not take the picture: %s", err.Error()))
		return
	}

	err = b.sendPic(msg, now, im)
	if err != nil {
		log.Err(err).Send()
		b.replyTo(msg, fmt.Sprintf("Could not send the picture: %s", err.Error()))
		return
	}
	log.Info().Interface("user", msg.From).Msg("sent picture")
}

func (b *Bot) RunForever() {
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	for update := range b.api.GetUpdatesChan(update) {
		log.Debug().Interface("update", update).Msg("new update")

		if update.Message == nil {
			log.Warn().Msg("no message, do not know how to handle this")
			continue
		}

		msg := update.Message
		if msg.From == nil || msg.Chat == nil {
			log.Warn().Msg("from or char is nil, discarding")
			continue
		}

		if !b.userAllowed(*msg.From) {
			log.Warn().Interface("from", msg.From).Msg("user not allowed")
			continue
		}

		switch msg.Text {
		case "/pic":
			b.handlePicCommand(msg)
		case "/start":
			b.replyTo(msg, "Use the /pic command to get a picture.")
		default:
			b.replyTo(msg, "I do not understand this. Use the /pic command to get a picture.")
		}
	}
}
