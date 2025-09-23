package minio

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	client *minio.Client
}

func NewMinIO(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*MinIO, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIO{client: client}, nil
}

func (m *MinIO) Upload(ctx context.Context, bucketName, objectKey string, data []byte, mimeType string) error {
	_, err := m.client.PutObject(ctx, bucketName, objectKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: mimeType})
	return err
}

func (m *MinIO) Download(ctx context.Context, bucketName, objectKey string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	// Read the object data into a buffer
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(obj); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *MinIO) Delete(ctx context.Context, bucketName, objectKey string) error {
	err := m.client.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinIO) EnsureBucket(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}
	return nil
}
