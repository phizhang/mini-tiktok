package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	Bucket string
	Client *s3.S3
}

func NewS3Storage(sess *session.Session, bucket string) *S3Storage {
	return &S3Storage{
		Bucket: bucket,
		Client: s3.New(sess),
	}
}

func (s *S3Storage) Upload(file multipart.File, filename string) (string, error) {
	log.Printf("[S3Storage] Uploading file %s to bucket %s", filename, s.Bucket)
	// Read file into memory (for demo; for large files, use streaming upload)
	buf := make([]byte, 0)
	chunk := make([]byte, 4096)
	for {
		n, err := file.Read(chunk)
		if n > 0 {
			buf = append(buf, chunk[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("[S3Storage] Error reading file: %v", err)
			return "", err
		}
	}
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Body:   aws.ReadSeekCloser(bytes.NewReader(buf)),
	}
	_, err := s.Client.PutObjectWithContext(context.Background(), input)
	if err != nil {
		log.Printf("[S3Storage] Failed to upload to S3: %v", err)
		return "", err
	}
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.Bucket, filename)
	log.Printf("[S3Storage] File uploaded successfully: %s", url)
	return url, nil
}
