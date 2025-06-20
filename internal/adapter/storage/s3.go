package storage

import (
	"context"
	"io"
	"path/filepath"
	"taskflow/internal/domain/repositories"
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

type fileRepository struct {
	s3Client *s3.S3
	bucket   string
}

func NewFileRepository(s3Client *s3.S3, bucket string) repositories.FileRepository {
	return &fileRepository{
		s3Client: s3Client,
		bucket:   bucket,
	}
}

func (r *fileRepository) Upload(ctx context.Context, filename string, content io.Reader, contentType string) (string, error) {
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

func (r *fileRepository) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	result, err := r.s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fileID),
	})

	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (r *fileRepository) Delete(ctx context.Context, fileID string) error {
	_, err := r.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fileID),
	})

	return err
}
