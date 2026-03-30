package telegram

import (
	"cyber-defender-bot-tg/internal/virustotal"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	api        *tgbotapi.BotAPI
	downloader *Downloader
	vtClient   *virustotal.Client
}

func NewHandler(
	api *tgbotapi.BotAPI,
	vtClient *virustotal.Client,
) *Handler {
	return &Handler{
		api:        api,
		downloader: NewDownloader(api),
		vtClient:   vtClient,
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("update received")

	message := update.Message

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

func (h *Handler) handleDocument(
	chatID int64,
	doc *tgbotapi.Document,
) {

	localPath, err := h.downloader.DownloadDocument(doc)
	if err != nil {
		log.Printf("download document: %v", err)

		h.sendMessage(
			chatID,
			"Не удалось скачать файл.",
		)

		return
	}

	defer os.Remove(localPath)

	log.Printf("file downloaded: %s", localPath)

	analysisID, err := h.vtClient.UploadFile(localPath)
	if err != nil {

		log.Printf("upload file: %v", err)

		h.sendMessage(
			chatID,
			"Не удалось отправить файл на проверку.",
		)

		return
	}

	log.Printf(
		"file uploaded to virustotal: %s",
		analysisID,
	)

	h.sendMessage(chatID, "Файл получен. Отправляю на проверку...")
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := h.api.Send(msg); err != nil {
		log.Printf("send message: %v", err)
	}
}
