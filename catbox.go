package catbox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"time"
)

// UploadFile uploads a file from a bytes.Buffer to CatBox.
func UploadFile(fileBuffer *bytes.Buffer, fileName string, timeout time.Duration, userHash string) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("fileToUpload", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(part, fileBuffer); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := writer.WriteField("reqtype", "fileupload"); err != nil {
		return "", fmt.Errorf("failed to write reqtype: %w", err)
	}

	if userHash != "" {
		if err := writer.WriteField("userHash", userHash); err != nil {
			return "", fmt.Errorf("failed to write userHash: %w", err)
		}
	}

	defer func() {
		if err := writer.Close(); err != nil {
			fmt.Printf("failed to close writer: %s\n", err)
		}
	}()

	req, err := http.NewRequest("POST", "https://catbox.moe/user/api.php", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() {
			return "", fmt.Errorf("upload request timed out after %d seconds", int(timeout.Seconds()))
		}
		return "", fmt.Errorf("failed to connect to Catbox: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP error occurred: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil
}
