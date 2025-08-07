package storage

import (
	"mime/multipart"
)

// Storage is the interface for video file storage backends
// Upload uploads a file and returns the storage URL or path, or an error
//
// file: the file to upload
// filename: the original filename
// returns: the URL or path where the file is stored, or error

type Storage interface {
	Upload(file multipart.File, filename string) (string, error)
}
