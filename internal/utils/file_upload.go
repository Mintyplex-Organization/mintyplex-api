package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
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

// func UploadToSia(file multipart.File, userID, bucket, filename string) (string, error) {
// 	siaServer := os.Getenv("SIA_SERVER")
// 	url := fmt.Sprintf("http://%s/api/worker/objects/users/%s/%s/%s", siaServer, userID, bucket, filename)
// 	fmt.Println(url)

// 	req, err := http.NewRequest("PUT", url, file)
// 	if err != nil {
// 		return "", err
// 	}
// 	req.Header.Set("Authorization", "Basic "+os.Getenv("SIA_API_AUTH"))
// 	req.Header.Set("Content-Type", "application/octet-stream")

// 	// Use a custom HTTP client to ignore SSL verification (not recommended for production)
// 	client := &http.Client{
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 		},
// 	}

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(res.Body)
// 		return "", fmt.Errorf("failed to upload to Sia renterd: %s", body)
// 	}

// 	return url, nil
// }

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

func UploadToSiaHandler(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error parsing form data: " + err.Error(),
		})
	}
	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No file uploaded",
		})
	}

	fileHead := files[0]
	file, err := fileHead.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to open file: " + err.Error(),
		})
	}
	defer file.Close()

	// Determine the file type based on the file extension
	fileType := DetermineFileType(fileHead.Filename)

	// Extract wallet address from user ID in the request parameters
	WalletAddress := c.Params("id")
	if WalletAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "User ID is required",
		})
	}

	// Set product and bucket manually for now
	bucket := fmt.Sprintf("%s-bucket", fileType)

	siaURL, err := UploadToSia(file, WalletAddress, bucket, fileHead.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to upload file to Sia renterd: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "File uploaded successfully",
		"url":     siaURL,
	})
}
