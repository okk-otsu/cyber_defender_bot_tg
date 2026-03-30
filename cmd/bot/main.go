package main

import (
	"cyber-defender-bot-tg/internal/config"
	"cyber-defender-bot-tg/internal/telegram"
	"cyber-defender-bot-tg/internal/virustotal"
)

func main() {
	cfg := config.MustLoad()

	vtClient := virustotal.NewClient(
		cfg.VirusTotalAPIKey,
	)

	bot := telegram.NewBot(
		cfg.TelegramBotToken,
		vtClient,
		&cfg,
	)

	bot.Run()
}
