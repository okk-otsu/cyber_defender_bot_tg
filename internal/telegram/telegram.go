package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(token string) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("create telegram bot api: %v", err)
	}

	return &Bot{api}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.api.GetUpdatesChan(u)

	log.Printf("bot started: @%s", b.api.Self.UserName)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		chatID := update.Message.Chat.ID

		switch text {
		case "/start":
			b.sendMessage(chatID, "Привет. Я бот для проверки файлов на угрозы.")
		case "/ping":
			b.sendMessage(chatID, "pong")
		default:
			b.sendMessage(chatID, "Я пока понимаю только /start и /ping")
		}
	}
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("send Message: %v", err)
	}
}
