package virustotal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	maxAnalysisAttempts  = 20
	analysisPollInterval = 3 * time.Second
	baseURL              = "https://www.virustotal.com/api/v3"
)

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
	if resp.StatusCode == http.StatusConflict {
		respBody, _ := io.ReadAll(resp.Body)
		return "", &AlreadySubmittedError{
			Message: string(respBody),
		}
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.Data.ID, nil
}

func (c *Client) GetAnalysis(analysisID string) (*AnalysisResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/analyses/"+analysisID,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-apikey", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result AnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) WaitForAnalysis(analysisID string) (*AnalysisResponse, error) {
	for i := 0; i < 20; i++ {
		result, err := c.GetAnalysis(analysisID)
		if err != nil {
			return nil, err
		}

		status := result.Data.Attributes.Status

		log.Printf("analysis status: %s, attempt: %d", status, i+1)

		if status == "completed" {
			return result, nil
		}

		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("analysis did not complete in time")
}
