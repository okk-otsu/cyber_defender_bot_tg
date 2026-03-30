package virustotal

import (
	"fmt"
	"net/http"
)

type Client struct {
	apiKey string
	client *http.Client
}

type UploadResponse struct {
	Data UploadResponseData `json:"data"`
}

type UploadResponseData struct {
	ID string `json:"id"`
}

type AnalysisResponse struct {
	Data AnalysisData `json:"data"`
}

type AnalysisData struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes AnalysisAttributes `json:"attributes"`
}

type AnalysisAttributes struct {
	Status string        `json:"status"`
	Stats  AnalysisStats `json:"stats"`
}

type AnalysisStats struct {
	Malicious  int `json:"malicious"`
	Suspicious int `json:"suspicious"`
	Harmless   int `json:"harmless"`
	Undetected int `json:"undetected"`
	Timeout    int `json:"timeout"`
}

type AlreadySubmittedError struct {
	Message string
}

func (e *AlreadySubmittedError) Error() string {
	return fmt.Sprintf("already submitted: %s", e.Message)
}
