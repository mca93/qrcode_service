package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/skip2/go-qrcode"
)

// generateQRCodeImage generates a QR code image and uploads it to an image server.
// It returns the image URL and the deep link.
func generateQRCodeImage(deepLink string, imageServerURL string) (string, string, error) {
	// Generate the QR code image
	qr, err := qrcode.New(deepLink, qrcode.High)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate QR code: %w", err)
	}
	qr.BackgroundColor = color.White
	qr.ForegroundColor = color.Black

	// Encode the QR code as PNG
	var buf bytes.Buffer
	if err := qr.Write(256, &buf); err != nil {
		return "", "", fmt.Errorf("failed to encode QR code: %w", err)
	}

	// Upload the image to the image server
	imageURL, err := uploadImageToServer(&buf, imageServerURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to upload QR code image: %w", err)
	}

	return imageURL, deepLink, nil
}

// uploadImageToServer uploads the QR code image to the image server and returns the image URL.
func uploadImageToServer(imageData io.Reader, imageServerURL string) (string, error) {
	// Create a new multipart form request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the image file to the form
	part, err := writer.CreateFormFile("file", "qrcode.png")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, imageData); err != nil {
		return "", fmt.Errorf("failed to copy image data: %w", err)
	}

	// Close the writer to finalize the form
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Send the POST request to the image server
	req, err := http.NewRequest("POST", imageServerURL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to upload image: received non-200 response")
	}

	// Parse the response to get the image URL
	var responseBody bytes.Buffer
	if _, err := io.Copy(&responseBody, resp.Body); err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Assuming the image server returns the image URL in the response body
	return responseBody.String(), nil
}
