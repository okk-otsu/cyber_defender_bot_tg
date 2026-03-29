package telegram

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Downloader struct {
	api *tgbotapi.BotAPI
}

func NewDownloader(api *tgbotapi.BotAPI) *Downloader {
	return &Downloader{
		api: api,
	}
}

func (d *Downloader) DownloadDocument(doc *tgbotapi.Document) (string, error) {
	if doc == nil {
		return "", fmt.Errorf("document is nil")
	}

	fileURL, err := d.buildFileURL(doc.FileID)
	if err != nil {
		return "", err
	}

	tmpFile, err := createTempFile(doc.FileName)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if err := downloadToFile(fileURL, tmpFile); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func (d *Downloader) buildFileURL(fileID string) (string, error) {
	fileConfig := tgbotapi.FileConfig{
		FileID: fileID,
	}

	tgFile, err := d.api.GetFile(fileConfig)
	if err != nil {
		return "", fmt.Errorf("get telegram file: %w", err)
	}

	return tgFile.Link(d.api.Token), nil
}

func createTempFile(fileName string) (*os.File, error) {
	file, err := os.CreateTemp(
		"",
		"tg-upload-*"+filepath.Ext(fileName),
	)
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}

	return file, nil
}

func downloadToFile(url string, file *os.File) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download telegram file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"download telegram file failed: status=%d",
			resp.StatusCode,
		)
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("save temp file: %w", err)
	}

	return nil
}
