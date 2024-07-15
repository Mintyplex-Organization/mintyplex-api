package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

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

func UploadToSia(file multipart.File, userID, bucket, filename string) (string, error) {
	uploadURL := "https://upload.mintyplex.com/s5/upload"
	authToken := os.Getenv("SIA_API_AUTH")
	fmt.Println(authToken)

	if authToken == "" {
		return "", fmt.Errorf("SIA_API_AUTH environment variable is not set")
	}

	// Create the request
	// Add necessary query parameters if required
	// For example, you might need to include bucket information in the URL
	url := fmt.Sprintf("%s?userID=%s&bucket=%s&filename=%s", uploadURL, userID, bucket, filename)
	fmt.Println("Request URL:", url)

	req, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		return "", err
	}

	// Set the Authorization header using the base64-encoded admin credentials
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/octet-stream")

	// Use a custom HTTP client to ignore SSL verification (not recommended for production)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("failed to upload to Sia: %s", body)
	}

	return url, nil
}
