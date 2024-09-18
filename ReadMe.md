# CatBox Go Package

The `catbox` package provides functions to upload files to CatBox and LitterBox, allowing for easy file sharing and storage through their APIs.

## Features

- **Upload files to CatBox**: Use `UploadFile` to upload a file to a specific user account.
- **Upload files to LitterBox**: Use `UploadToLitterBox` to upload files with a specified storage duration.

### ToDo's
- [ ] deleteFiles(files, userHash) - Delete files from CatBox with a specific user account.
- [ ] uploadAlbum(files, title, description, userHash): Create an album with the uploaded files.
- [ ] deleteAlbum, editAlbum, createAlbum ... etc

## Installation

To use the `catbox` package in your Go project, you can import it directly:

```go
import "github.com/AshokShau/catbox"
```

## Usage
Examples of uploading files to CatBox and LitterBox are shown below.

RealTime Example: [here](example/main.go)


### Uploading a File to CatBox

```go
fileBuffer := bytes.NewBuffer(yourFileBytes)
fileName := "example.txt"
timeout := 10 * time.Second
userHash := "" // optional

url, err := catbox.UploadFile(fileBuffer, fileName, timeout, userHash)
if err != nil {
    log.Fatalf("Upload failed: %v", err)
}
fmt.Println("Uploaded to CatBox:", url)
```

### Uploading a File to LitterBox

```go
fileBuffer := bytes.NewBuffer(yourFileBytes)
fileName := "example.txt"
duration := "24h" // Options: '1h', '12h', '24h', '72h', '1w'
timeout := 10 * time.Second

url, err := catbox.UploadToLitterBox(fileBuffer, fileName, duration, timeout)
if err != nil {
    log.Fatalf("Upload failed: %v", err)
}       
fmt.Println("Uploaded to LitterBox:", url)
```

## Parameters

### `UploadFile`

- `fileBuffer`: A `*bytes.Buffer` containing the file content.
- `fileName`: The name of the file to be uploaded.
- `timeout`: Duration to wait for the server response.
- `userHash`: Optional parameter to upload files to a specific user account.

### `UploadToLitterBox`

- `fileBuffer`: A `*bytes.Buffer` containing the file content.
- `fileName`: The name of the file to be uploaded.
- `duration`: Duration for which the file should be stored on the server (e.g., '1h', '12h', '24h', '72h', '1w').
- `timeout`: Duration to wait for the server response.

## Error Handling

Both functions return an error if the upload fails. Common reasons include:

- Network issues
- Invalid parameters
- Non-200 HTTP responses

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [CatBox](https://catbox.moe/) for their file hosting service.
- [LitterBox](https://litterbox.catbox.moe/) for their file storage service.
- [Golang](https://golang.org/) for the Go programming language.