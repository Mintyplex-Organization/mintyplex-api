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
func CreateBucket(bucketName string) error {
	siaServer := os.Getenv("SIA_SERVER")
	authToken := os.Getenv("SIA_API_AUTH")
	if authToken == "" {
		return fmt.Errorf("SIA_API_AUTH environment variable is not set")
	}

	url := fmt.Sprintf("http://%s/api/worker/buckets/%s", siaServer, bucketName)
	fmt.Println("Creating bucket at URL:", url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to create bucket: %s", body)
	}

	fmt.Println("Bucket created successfully")
	return nil
}

func UploadToSia(file multipart.File, userID, bucket, filename string) (string, error) {
	uploadURL := "https://upload.mintyplex.com/s5/upload"
	// uploadURL := "localhost:9980/s5/upload"
	authToken := os.Getenv("SIA_API_AUTH")
	fmt.Println(authToken)

	if authToken == "" {
		return "", fmt.Errorf("SIA_AUTH_TOKEN environment variable is not set")
	}

	// Create the request
	// url := fmt.Sprintf("%s?auth_token=%s", uploadURL, authToken)
	// fmt.Println("Request URL:", url)

	url := fmt.Sprintf("%s/%s", uploadURL, bucket)
	fmt.Println("Uploading file to URL:", url)

	req, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		return "", err
	}
	// Optionally, you can set the Authorization header instead of using the query parameter
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SIA_API_AUTH"))

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

	return uploadURL, nil
}
