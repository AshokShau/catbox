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
// The fileName parameter specifies the name of the file to be uploaded.
// The timeout parameter specifies how long the client should wait for the server to respond.
// The userHash parameter is optional and can be used to upload files to a specific user account.
// The function returns the URL of the uploaded file or an error if the upload failed.
func UploadFile(fileBuffer *bytes.Buffer, fileName string, timeout time.Duration, userHash string) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if part, err := writer.CreateFormFile("fileToUpload", fileName); err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	} else if _, err = io.Copy(part, fileBuffer); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	_ = writer.WriteField("reqtype", "fileupload")
	if userHash != "" {
		_ = writer.WriteField("userHash", userHash)
	}
	_ = writer.Close()

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

// UploadToLitterBox uploads a file from a bytes.Buffer to LitterBox.
// The fileName parameter specifies the name of the file to be uploaded.
// The duration parameter specifies how long the file should be stored on the server, Options: '1h', '12h', '24h', '72h', '1w'.
// The timeout parameter specifies how long the client should wait for the server to respond.
func UploadToLitterBox(fileBuffer *bytes.Buffer, fileName string, duration string, timeout time.Duration) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if part, err := writer.CreateFormFile("fileToUpload", fileName); err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	} else if _, err = io.Copy(part, fileBuffer); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	_ = writer.WriteField("reqtype", "fileupload")
	_ = writer.WriteField("time", duration)
	_ = writer.Close()

	req, err := http.NewRequest("POST", "https://litterbox.catbox.moe/resources/internals/api.php", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() {
			return "", fmt.Errorf("upload to Litterbox timed out after %d seconds", int(timeout.Seconds()))
		}
		return "", fmt.Errorf("failed to connect to Litterbox: %w", err)
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
