package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

const (
	MaxImageSize = 5 << 20 // 5MB
	UploadDir = "uploads/images"
)

var (
	ErrFileTooLarge = errors.New("file exceeds maximum allowed size of 5MB")
	ErrInvalidFileType = errors.New("unsupported file type: only JPEG, PNG, and GIF")
	ErrEmptyFile = errors.New("file is empty")
)

var magicBytes = map[string][]byte{
	"image/jpeg": {0xFF, 0xD8, 0xFF},
	"image/png": {0x89, 0x50, 0x4E, 0x47},
	"image/gif": {0x47, 0x49, 0x46, 0x38},
}

var extensions = map[string]string {
	"image/jpeg": ".jpg",
	"image/png": ".png",
	"image/gif": ".gif",
}

func detectMIMEType(file multipart.File) (string, error) {
	header := make([]byte, 8)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	if n == 0 {
		return "", ErrEmptyFile
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	for mime, magic := range magicBytes {
		if len(header) >= len(magic) && matchesMagic(header, magic) {
			return mime, nil
		}
	}
	return "", ErrInvalidFileType
}

func matchesMagic(header, magic []byte) bool {
	for i, b := range magic {
		if header[i] != b {
			return false
		}
	}
	return true
}

func SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size == 0 {
		return "", ErrEmptyFile
	}

	if header.Size > MaxImageSize {
		return "", ErrFileTooLarge
	}

	mimeType, err := detectMIMEType(file)
	if err != nil {
		return "", err
	}

	ext := extensions[mimeType]

	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	filename := id.String() + ext
	filePath := filepath.Join(UploadDir, filename)

	destination, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		// clean up partial file on failure
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}

// remove an image file from the filesystem
func DeleteImage(filePath string) error {
	if filePath == "" {
		return nil 
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}
