package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (s *MinioService) UploadFile(file io.Reader, size int64, originalFilename string) (string, error) {
	extension := filepath.Ext(originalFilename)
	objectName := uuid.New().String() + extension

	_, err := s.Client.PutObject(context.Background(), s.Bucket, objectName, file, size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return "http://" + s.Endpoint + "/" + s.Bucket + "/" + objectName, nil
}

func (s *MinioService) DeleteFile(imageURL string) error {
	expectedPrefix := "http://" + s.Endpoint + "/" + s.Bucket + "/"
	if !strings.HasPrefix(imageURL, expectedPrefix) {
		return fmt.Errorf("invalid image URL: unexpected origin")
	}
	objectName := filepath.Base(imageURL)
	return s.Client.RemoveObject(context.Background(), s.Bucket, objectName, minio.RemoveObjectOptions{})
}

func NewMinioService(endpoint string, accessKey string, secretKey string, bucket string) (*MinioService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return &MinioService{}, err
	}

	exists, err := client.BucketExists(context.Background(), bucket)
	if err != nil{
		return &MinioService{}, err
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
		if err != nil {
			return &MinioService{}, err
		}
	}

	return &MinioService{ Client: client, Bucket: bucket, Endpoint: endpoint}, nil
}