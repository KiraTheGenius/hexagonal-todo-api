package storage

import (
	"context"
	"io"
	"path/filepath"
	"taskflow/internal/domain/shared"
	"taskflow/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

func NewS3Client(cfg config.S3Config) *s3.S3 {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
	}))

	return s3.New(sess)
}

type s3Storage struct {
	s3Client *s3.S3
	bucket   string
}

func NewS3Storage(s3Client *s3.S3, bucket string) shared.Storage {
	return &s3Storage{
		s3Client: s3Client,
		bucket:   bucket,
	}
}

func (r *s3Storage) Upload(ctx context.Context, filename string, content io.Reader, contentType string) (string, error) {
	fileID := uuid.New().String() + filepath.Ext(filename)

	_, err := r.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(fileID),
		Body:        aws.ReadSeekCloser(content),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	return fileID, nil
}

func (r *s3Storage) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	result, err := r.s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fileID),
	})

	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (r *s3Storage) Delete(ctx context.Context, fileID string) error {
	_, err := r.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fileID),
	})

	return err
}

func (r *s3Storage) GetURL(ctx context.Context, fileID string) (string, error) {
	// Generate a presigned URL for the file
	req, _ := r.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fileID),
	})

	url, err := req.Presign(15 * 60) // 15 minutes
	if err != nil {
		return "", err
	}

	return url, nil
}
