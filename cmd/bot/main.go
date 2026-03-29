package main

import (
	"cyber-defender-bot-tg/internal/config"
	"cyber-defender-bot-tg/internal/telegram"
)

func main() {
	cfg := config.MustLoad()

	bot := telegram.NewBot(cfg.TelegramBotToken)
	bot.Run()
}
