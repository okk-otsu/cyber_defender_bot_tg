package telegram

import (
	"cyber-defender-bot-tg/internal/config"
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
	cfg *config.Config,
) *Bot {

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("create telegram bot:", err)
	}

	handler := NewHandler(api, vtClient, cfg)

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
