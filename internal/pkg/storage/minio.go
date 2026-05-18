package storage

import (
	"context"
	"fmt"

	"github.com/SovetkanB/FlipFlow/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	client *minio.Client
	cfg    config.MinIOConfig
}

func NewMinIO(cfg config.MinIOConfig) (*MinIO, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("minio.New: %w", err)
	}

	return &MinIO{client: client, cfg: cfg}, nil
}

func (m *MinIO) EnsureBuckets(ctx context.Context) error {
	for _, bucket := range []string{m.cfg.BucketPhotos, m.cfg.BucketFiles} {
		exists, err := m.client.BucketExists(ctx, bucket)
		if err != nil {
			return err
		}

		if !exists {
			if err := m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return fmt.Errorf("make bucket %s: %w", bucket, err)
			}
		}
	}

	return nil
}
