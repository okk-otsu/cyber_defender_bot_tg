package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	VirusTotalAPIKey string
	MaxFileSizeBytes int64
}

func MustLoad() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment")
	}

	return Config{
		TelegramBotToken: mustEnv("TELEGRAM_BOT_TOKEN"),
		VirusTotalAPIKey: mustEnv("VIRUSTOTAL_API_KEY"),
		MaxFileSizeBytes: mustEnvInt64("MAX_FILE_SIZE_BYTES"),
	}
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not set", key)
	}

	return value
}

func mustEnvInt64(key string) int64 {
	value := mustEnv(key)

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("%s is not a number", key)
	}
	return n
}
