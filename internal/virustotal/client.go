package virustotal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const baseURL = "https://www.virustotal.com/api/v3"

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) UploadFile(filePath string) (string, error) {
	body, contentType, err := buildMultipartBody(filePath)
	if err != nil {
		return "", err
	}

	req, err := c.newUploadRequest(body, contentType)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	return decodeUploadResponse(resp)
}

func buildMultipartBody(filePath string) (*bytes.Buffer, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, "", fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, "", fmt.Errorf("copy file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close writer: %w", err)
	}

	return body, writer.FormDataContentType(), nil
}

func (c *Client) newUploadRequest(body *bytes.Buffer, contentType string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, baseURL+"/files", body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-apikey", c.apiKey)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

func decodeUploadResponse(resp *http.Response) (string, error) {
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result UploadResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Data.ID, nil
}
