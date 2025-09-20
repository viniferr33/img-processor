package image

import "context"

type Image struct {
	ID           string
	OriginalName string
	Size         int64
	MimeType     string
	Width        int
	Height       int
	ParentID     string
	Version      int
	CreatedAt    int64
	UpdatedAt    int64

	BucketName string
	ObjectKey  string
}

type ImageRepository interface {
	GetByID(ctx context.Context, id string) (*Image, error)
	Create(ctx context.Context, image *Image) error
	Update(ctx context.Context, image *Image) error
	Delete(ctx context.Context, id string) error
}

type ObjectStorage interface {
	Upload(ctx context.Context, bucketName, objectKey string, data []byte, mimeType string) error
	Download(ctx context.Context, bucketName, objectKey string) ([]byte, error)
	Delete(ctx context.Context, bucketName, objectKey string) error
}

func NewImage(id, originalName string, size int64, mimeType string, width, height int, parentID string, version int, bucketName, objectKey string) *Image {
	return &Image{
		ID:           id,
		OriginalName: originalName,
		Size:         size,
		MimeType:     mimeType,
		Width:        width,
		Height:       height,
		ParentID:     parentID,
		Version:      version,
		BucketName:   bucketName,
		ObjectKey:    objectKey,
	}
}
