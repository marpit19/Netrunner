package http

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/appyzdl/Netrunner/pkg/http/status"
)

func getContentType(path string) string {
	ext := filepath.Ext(path)

	// First, try to use the standard library's mime.TypeByExtension
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}

	// If the standard library doesn't recognize the extension, use our own mapping
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".xml":
		return "application/xml"
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".mp3":
		return "audio/mpeg"
	case ".mp4":
		return "video/mp4"
	default:
		return "application/octet-stream" // Default to binary data
	}
}

func StaticFileHandler(basePath string) HandlerFunc {
	return func(req *Request) *Response {
		// Remove the "/static" prefix from the request path
		filePath := strings.TrimPrefix(req.Path, "/static")

		// If the path is empty, serve index.html
		if filePath == "" || filePath == "/" {
			filePath = "/index.html"
		}

		// Ensure the path doesn't try to access parent directories
		if strings.Contains(filePath, "..") {
			return NotFoundResponse()
		}

		fullPath := filepath.Join(basePath, filePath)

		fmt.Printf("Attempting to serve file: %s\n", fullPath) // Debug log

		file, err := os.Open(fullPath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err) // Debug log
			return NotFoundResponse()
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			fmt.Printf("Error getting file stats: %v\n", err) // Debug log
			return InternalServerErrorResponse()
		}

		if stat.IsDir() {
			fmt.Println("Requested path is a directory") // Debug log
			return NotFoundResponse()
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err) // Debug log
			return InternalServerErrorResponse()
		}

		resp := NewResponse()
		resp.StatusCode = status.OK
		resp.StatusText = StatusText(status.OK)
		resp.SetHeader("Content-Type", getContentType(fullPath))
		resp.SetHeader("Content-Length", fmt.Sprintf("%d", len(content)))
		resp.Body = content
		return resp
	}
}
