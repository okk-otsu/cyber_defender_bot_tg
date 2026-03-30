package virustotal

import "net/http"

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
