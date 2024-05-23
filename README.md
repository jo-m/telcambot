# Telegram Webcam Bot

Telegram bot which sends pictures from a webcam upon request.

A single command is supported: `/pic`.

Uses [libcamera-jpeg](https://www.raspberrypi.com/documentation/computers/camera_software.html#libcamera-jpeg) and is tested only
on a Raspberry Pi Zero W.

Written in Go, thanks go out to <https://github.com/go-telegram-bot-api/telegram-bot-api>.

## Usage

1. Set up a Telegram bot: <https://core.telegram.org/bots>
2. Set up the bot with the commands description `pic - request a picture` (via BotFather)
3. Copy `.env.example` to `.env` and add your bot token.

```
$ go run ./ --help
Usage: telecambot [--log-pretty] [--log-level LEVEL] --bot-api-token TOK [--bot-debug] [--bot-users IDS,] [--pic-cache-dur DUR] [--libcamera-camera-ix N] [--libcamera-w W] [--libcamera-h H]

Options:
  --log-pretty           log pretty [default: true, env: LOG_PRETTY]
  --log-level LEVEL      log level [default: info, env: LOG_LEVEL]
  --bot-api-token TOK    get it from https://t.me/Botfather [env: BOT_API_TOKEN]
  --bot-debug            run telegram bot in debug mode [default: false, env: BOT_DEBUG]
  --bot-users IDS,       whitelist of Telegram user ids allowed to use the bot [env: BOT_USERS]
  --pic-cache-dur DUR    for how long to cache pictures [default: 10s, env: PIC_CACHE_DUR]
  --libcamera-camera-ix N
                         libcamera camera number, see 'libcamera-jpeg --list-cameras' [default: 0, env: LIBCAMERA_CAMERA_IX]
  --libcamera-w W        width [default: 1920, env: LIBCAMERA_W]
  --libcamera-h H        height [default: 1080, env: LIBCAMERA_H]
  --help, -h             display this help and exit
```
