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
	CID      string `json:"cid"`
	FileName string `json:"filename"`
	FileType string `json:"filetype"`
	Size     int64  `json:"size"`
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

func UploadToSia(file multipart.File, fileSize int64, userID, bucket, filename string) (SiaUploadResponse, error) {
	// var siaResp SiaUploadResponse

	uploadURL := "https://upload.mintyplex.com/s5/upload/"
	authToken := os.Getenv("SIA_API_AUTH")

	if authToken == "" {
		fmt.Println("missing variable")
		return SiaUploadResponse{}, fmt.Errorf("SIA_AUTH_TOKEN environment variable is not set")
	}

	req, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		fmt.Println("bad request")
		return SiaUploadResponse{}, err
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
		return SiaUploadResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Printf("Error response body: %s\n", body) // Log the error body for debugging
		return SiaUploadResponse{}, fmt.Errorf("failed to upload to Sia: %s", body)
	}

	var siaResponse SiaUploadResponse
	if err := json.NewDecoder(res.Body).Decode(&siaResponse); err != nil {
		return SiaUploadResponse{}, fmt.Errorf("fialed to parse response: %s",err, err.Error())
	}

	siaResponse.FileName = filename
	siaResponse.FileType = DetermineFileType(filename)
	siaResponse.Size = fileSize

	return siaResponse, nil
}
