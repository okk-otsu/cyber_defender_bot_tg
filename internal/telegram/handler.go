package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	api        *tgbotapi.BotAPI
	downloader *Downloader
}

func NewHandler(api *tgbotapi.BotAPI) *Handler {
	return &Handler{
		api:        api,
		downloader: NewDownloader(api),
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("update received")

	message := update.Message
	log.Printf("message received")

	if message.Document != nil {
		log.Printf("document received: %s", message.Document.FileName)
	}

	if len(message.Photo) > 0 {
		log.Printf("photo received")
	}

	chatID := message.Chat.ID

	switch {
	case message.Document != nil:
		h.handleDocument(chatID, message.Document)

	case message.Text == "/start":
		h.sendMessage(chatID, "Привет. Отправь мне файл как document.")

	case message.Text == "/ping":
		h.sendMessage(chatID, "pong")

	default:
		h.sendMessage(chatID, "Я пока понимаю /start, /ping и document.")
	}
}

func (h *Handler) handleDocument(chatID int64, doc *tgbotapi.Document) {
	localPath, err := h.downloader.DownloadDocument(doc)
	if err != nil {
		log.Printf("download document: %v", err)
		h.sendMessage(chatID, "Не удалось скачать файл.")
		return
	}

	h.sendMessage(
		chatID,
		fmt.Sprintf(
			"Файл скачан.\nИмя: %s\nРазмер: %d байт",
			doc.FileName,
			doc.FileSize,
		),
	)
	log.Printf("file downloaded: %s", localPath)
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := h.api.Send(msg); err != nil {
		log.Printf("send message: %v", err)
	}
}
