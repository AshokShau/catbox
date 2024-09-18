package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/AshokShau/catbox"
)

func main() {
	// Open the image file
	filePath := "example/img.png"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Read the file into a buffer
	var fileBuffer bytes.Buffer
	if _, err = fileBuffer.ReadFrom(file); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Set the file name and user hash (if any)
	fileName := "image.png"
	userHash := "" // Optional: Set your user hash if you have one

	// Set the timeout duration
	timeout := 10 * time.Second

	// Upload the file to CatBox
	response, err := catbox.UploadFile(&fileBuffer, fileName, timeout, userHash)
	if err != nil {
		fmt.Printf("Upload to CatBox failed: %v\n", err)
		return
	}
	fmt.Println("Upload to CatBox successful! Response:", response)

	// Re-open the file for the second upload
	file, err = os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Reset the file buffer and read the file again
	fileBuffer.Reset()
	if _, err = fileBuffer.ReadFrom(file); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Upload the file to LitterBox
	duration := "1h" // Duration for which the file should be stored
	response, err = catbox.UploadToLitterBox(&fileBuffer, fileName, duration, timeout)
	if err != nil {
		fmt.Printf("Upload to LitterBox failed: %v\n", err)
		return
	}
	fmt.Println("Upload to LitterBox successful! Response:", response)
}
