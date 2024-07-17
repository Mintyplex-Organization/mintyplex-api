package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type SiaUploadResponse struct {
	CID string `json:"cid"`
	// Add other fields if necessary
}

func DetermineFileType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".mp3"), strings.HasSuffix(filename, ".m4a"), strings.HasSuffix(filename, ".wav"):
		return "audio"
	case strings.HasSuffix(filename, ".jpg"), strings.HasSuffix(filename, ".jpeg"), strings.HasSuffix(filename, ".webp"), strings.HasSuffix(filename, ".png"), strings.HasSuffix(filename, ".gif"):
		return "image"
	case strings.HasSuffix(filename, ".pdf"), strings.HasSuffix(filename, ".epub"), strings.HasSuffix(filename, ".txt"), strings.HasSuffix(filename, ".mobi"):
		return "ebook"
	default:
		return "unknown"
	}
}

func UploadToSia(file multipart.File, userID, bucket, filename string) (SiaUploadResponse, error) {
	var siaResp SiaUploadResponse

	uploadURL := "https://upload.mintyplex.com/s5/upload/"
	authToken := os.Getenv("SIA_API_AUTH")

	if authToken == "" {
		return siaResp, fmt.Errorf("SIA_AUTH_TOKEN environment variable is not set")
	}

	req, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		return siaResp, err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return siaResp, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return siaResp, fmt.Errorf("failed to upload to Sia: %s", body)
	}

	// Decode the JSON response
	err = json.NewDecoder(res.Body).Decode(&siaResp)
	if err != nil {
		return siaResp, fmt.Errorf("failed to decode Sia response: %w", err)
	}

	return siaResp, nil
}
