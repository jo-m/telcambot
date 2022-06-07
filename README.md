# Telegram Webcam Bot

Telegram bot which sends pictures from a webcam upon request.

A single command is supported: `/pic`.

Works only on Linux (uses the video4linux API).

Written in Go, thanks go out to <https://github.com/blackjack/webcam> and <https://github.com/go-telegram-bot-api/telegram-bot-api>.

## Usage

1. Set up a Telegram bot: <https://core.telegram.org/bots>
2. Set up the bot with the commands description `pic - request a picture` (via BotFather)
3. Copy `.env.example` to `.env` and add your bot token.

```
$ go run ./cmd/telcambot/ --help
Usage: telcambot [--log-pretty] [--log-level LEVEL] --bot-api-token TOK [--bot-debug] [--bot-users IDS,] [--pic-cache-dur DUR] [--pic-dev DEV] [--pic-format STR] [--pic-timeout DUR] [--pic-skip-frames N]

Options:
  --log-pretty           log pretty [default: false, env: LOG_PRETTY]
  --log-level LEVEL      log level [default: info, env: LOG_LEVEL]
  --bot-api-token TOK    get it from https://t.me/Botfather [env: BOT_API_TOKEN]
  --bot-debug            run telegram bot in debug mode [default: false, env: BOT_DEBUG]
  --bot-users IDS,       whitelist of Telegram user ids allowed to use the bot [env: BOT_USERS]
  --pic-cache-dur DUR    for how long to cache pictures [default: 10s, env: PIC_CACHE_DUR]
  --pic-dev DEV          camera video device file path [default: /dev/video3, env: PIC_DEV]
  --pic-format STR       camera preferred image format [default: Motion-JPEG, env: PIC_FORMAT]
  --pic-timeout DUR      how long to give camera time to start [default: 2s, env: PIC_TIMEOUT]
  --pic-skip-frames N    how many frames to skip until picture snap [default: 15, env: PIC_SKIP_FRAMES]
  --help, -h             display this help and exit
```

## TODOs

- [ ] Log to syslog
