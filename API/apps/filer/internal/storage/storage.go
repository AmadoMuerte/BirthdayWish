package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/filer/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	Client     *minio.Client
	BucketName string
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.DB.Host, cfg.DB.Port)
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.DB.User,
			cfg.DB.Pass,
			"",
		),
		Secure: cfg.DB.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DB.Timeout)
	defer cancel()

	if _, err = client.ListBuckets(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	storage := &Storage{
		Client:     client,
		BucketName: cfg.DB.BucketName,
	}

	exists, err := client.BucketExists(ctx, cfg.DB.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		if err = client.MakeBucket(ctx, cfg.DB.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		if err := storage.SetBucketPolicy(); err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	return storage, nil
}

func (s *Storage) SetBucketPolicy() error {
	policy := `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::` + s.BucketName + `/*"
        },
        {
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:ListBucket",
            "Resource": "arn:aws:s3:::` + s.BucketName + `"
        }
    ]
}`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.Client.SetBucketPolicy(ctx, s.BucketName, policy)
}
