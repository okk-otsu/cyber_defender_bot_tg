package main

import (
	"fmt"

	"cyber-defender-bot-tg/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("telegram token loaded:", cfg.TelegramBotToken != "")
	fmt.Println("virustotal api key loaded:", cfg.VirusTotalAPIKey != "")
	fmt.Println("max file size bytes:", cfg.MaxFileSizeBytes)
}
