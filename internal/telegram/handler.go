package telegram

import (
	"cyber-defender-bot-tg/internal/virustotal"
	"errors"
	"fmt"
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
		h.sendMessage(chatID, startMessageText())

	case message.Text == "/help":
		h.sendMessage(chatID, helpMessageText())

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
		h.sendMessage(chatID, "Не удалось скачать файл.")
		return
	}
	defer os.Remove(localPath)

	log.Printf("file downloaded: %s", localPath)

	h.sendMessage(chatID, "Файл получен. Отправляю на проверку...")

	analysisID, err := h.vtClient.UploadFile(localPath)
	if err != nil {
		var alreadySubmittedErr *virustotal.AlreadySubmittedError
		if errors.As(err, &alreadySubmittedErr) {
			log.Printf("file already submitted: %v", err)
			h.sendMessage(chatID, "Этот файл уже отправлен на проверку. Подожди немного и попробуй снова.")
			return
		}

		log.Printf("upload file: %v", err)
		h.sendMessage(chatID, "Не удалось отправить файл на проверку.")
		return
	}

	log.Printf("file uploaded to virustotal: %s", analysisID)

	result, err := h.vtClient.WaitForAnalysis(analysisID)
	if err != nil {
		log.Printf("wait for analysis: %v", err)
		h.sendMessage(chatID, "Файл отправлен на проверку, но не удалось получить результат.")
		return
	}

	stats := result.Data.Attributes.Stats

	text := buildVerdictText(doc.FileName, stats)

	h.sendMessage(chatID, text)
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := h.api.Send(msg); err != nil {
		log.Printf("send message: %v", err)
	}
}

func buildVerdictText(
	filename string,
	stats virustotal.AnalysisStats,
) string {

	totalEngines :=
		stats.Malicious +
			stats.Suspicious +
			stats.Harmless +
			stats.Undetected

	var verdict string

	switch {
	case stats.Malicious > 0:
		verdict = "⚠️ Обнаружены вредоносные срабатывания."

	case stats.Suspicious > 0:
		verdict = "⚠️ Обнаружены подозрительные срабатывания."

	default:
		verdict = "✅ Угроз не обнаружено."
	}

	return fmt.Sprintf(
		"Файл: %s\n\n"+
			"%s\n\n"+
			"Проверено антивирусами: %d\n\n"+
			"Результаты:\n"+
			"Вредоносные срабатывания: %d\n"+
			"Подозрительные срабатывания: %d\n"+
			"Без обнаруженных угроз: %d\n",

		filename,
		verdict,
		totalEngines,
		stats.Malicious,
		stats.Suspicious,
		stats.Undetected,
	)
}

func startMessageText() string {
	return `👋 Привет!

Я — бот для проверки файлов на вирусы через сервис VirusTotal.

Что я умею:
📄 Проверяю файлы на наличие вредоносного кода
🛡 Использую несколько антивирусных движков
⏳ Проверка обычно занимает до 1 минуты

Как пользоваться:
1. Отправь файл как "документ".
2. Подожди результат проверки.

ℹ️ Для списка команд напиши /help`
}

func helpMessageText() string {
	return `📖 Справка по использованию

Доступные команды:

/start — показать приветствие
/help — показать эту справку
/ping — проверить, работает ли бот

Как проверить файл:

1. Отправь файл в чат как "документ"
2. Бот загрузит файл
3. Отправит его на проверку
4. Вернёт результат

Ограничения:

📦 Максимальный размер файла: 10 МБ
⏳ Проверка может занять до 1 минуты

Важно:

⚠️ Даже если угрозы не обнаружены,
это не даёт 100% гарантии безопасности.`
}
