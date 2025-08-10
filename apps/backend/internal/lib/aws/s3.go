package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sriniously/tasker/internal/server"
)

type S3Client struct {
	server *server.Server
	client *s3.Client
}

func NewS3Client(server *server.Server, cfg aws.Config) *S3Client {
	return &S3Client{
		server: server,
		client: s3.NewFromConfig(cfg),
	}
}

func (s *S3Client) UploadFile(ctx context.Context, bucket string, fileName string, file io.Reader) (string, error) {
	fileKey := fmt.Sprintf("%s_%d", fileName, time.Now().Unix())

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileKey),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(http.DetectContentType(buffer.Bytes())),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return fileKey, nil
}

func (s *S3Client) CreatePresignedUrl(ctx context.Context, bucket string, objectKey string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	expiration := time.Minute * 60

	presignedUrl, err := presignClient.PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(expiration))
	if err != nil {
		return "", err
	}

	return presignedUrl.URL, nil
}

func (s *S3Client) DeleteObject(ctx context.Context, bucket string, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}

	return nil
}
