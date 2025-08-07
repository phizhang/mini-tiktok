package storage

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalStorage implements Storage interface, stores files on local disk

type LocalStorage struct {
	BasePath string // directory to store files
}

func (l *LocalStorage) Upload(file multipart.File, filename string) (string, error) {
	path := filepath.Join(l.BasePath, filename)
	log.Printf("[LocalStorage] Saving file to %s", path)

	if err := os.MkdirAll(l.BasePath, os.ModePerm); err != nil {
		log.Printf("[LocalStorage] Failed to create directory: %v", err)
		return "", err
	}
	out, err := os.Create(path)
	if err != nil {
		log.Printf("[LocalStorage] Failed to create file: %v", err)
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("[LocalStorage] Failed to copy file: %v", err)
		return "", err
	}
	log.Printf("[LocalStorage] File saved successfully: %s", path)
	return path, nil
}
