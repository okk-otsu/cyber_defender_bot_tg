package telegram

import (
	"cyber-defender-bot-tg/internal/virustotal"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler *Handler
}

func NewBot(
	token string,
	vtClient *virustotal.Client,
) *Bot {

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("create telegram bot api: %v", err)
	}

	handler := NewHandler(api, vtClient)

	return &Bot{
		api:     api,
		handler: handler,
	}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.api.GetUpdatesChan(u)

	log.Printf("bot started: @%s", b.api.Self.UserName)

	for update := range updates {
		b.handler.HandleUpdate(update)
	}
}
