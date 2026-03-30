# Cyber Defender Telegram Bot

Telegram-бот для проверки файлов на вредоносные программы с использованием VirusTotal.

Бот принимает файл от пользователя, отправляет его на проверку в VirusTotal, ожидает завершения анализа и возвращает результат проверки.

---

## Возможности

- Проверка файлов через VirusTotal API
- Поддержка различных типов файлов (pdf, exe, zip, изображения и др.)
- Ограничение размера файла через переменные окружения
- Автоматическая загрузка файлов из Telegram
- Ожидание результата анализа
- Возврат понятного пользователю результата
- Поддержка команд:
  - `/start`
  - `/help`
  - `/ping`

---

## Как работает бот

1. Пользователь отправляет файл **как документ**
2. Бот скачивает файл
3. Бот отправляет файл в VirusTotal
4. Бот ожидает завершения анализа
5. Бот возвращает результат проверки

---

## Требования

- Go **1.25+**
- Telegram Bot Token
- VirusTotal API Key

---

## Установка и запуск

### 1. Клонировать репозиторий

```bash
git clone https://github.com/okk-otsu/cyber_defender_bot_tg.git
cd cyber_defender_bot_tg
```

### 2. Создать файл `.env`:
```env
TELEGRAM_BOT_TOKEN=your_telegram_token
VIRUSTOTAL_API_KEY=your_virustotal_api_key
MAX_FILE_SIZE_BYTES=33554432
```

### 3. Установить зависимости:
```bash
go mod tidy
```

### 4. Запустить бота:
```bash
go run ./cmd/bot
```
---

## Структура проекта:
```
cyber-defender-bot-tg/
│
├── cmd/
│   └── bot/
│       └── main.go          # точка входа в приложение
│
├── internal/
│   ├── config/
│   │   ├── config.go        # загрузка конфигурации из .env
│   │   └── doc.go
│   │
│   ├── telegram/
│   │   ├── bot.go           # запуск Telegram-бота
│   │   ├── doc.go           
│   │   ├── handler.go       # обработка сообщений
│   │   └── downloader.go    # загрузка файлов из Telegram
│   │
│   └── virustotal/
│       ├── client.go        # клиент для работы с VirusTotal API
│       ├── doc.go
│       └── types.go         # структуры ответов API
│
├── .env.example             # пример переменных окружения
├── .gitignore
├── README.md
├── go.sum
└── go.mod
```
