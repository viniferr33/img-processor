package image

import (
	"context"
	"encoding/base64"
	"image"
	"strings"

	"github.com/google/uuid"
)

type ImageService struct {
	repo          ImageRepository
	storage       ObjectStorage
	defaultBucket string
}

func NewImageService(repo ImageRepository, storage ObjectStorage, defaultBucket string) *ImageService {
	return &ImageService{
		repo:          repo,
		storage:       storage,
		defaultBucket: defaultBucket,
	}
}

func (s *ImageService) UploadImage(ctx context.Context, data []byte, filename, mimetype, parentId, ownerId string, version int) (*Image, error) {
	imgWidth, imgHeight, err := GetImageMetadata(data)
	if err != nil {
		return nil, err
	}

	image := NewImage(
		uuid.New().String(),
		filename,
		int64(len(data)),
		mimetype,
		imgWidth,
		imgHeight,
		parentId,
		version,
		s.defaultBucket,
		uuid.New().String(),
		ownerId,
	)

	if err := s.storage.Upload(ctx, s.defaultBucket, image.ObjectKey, data, mimetype); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, image); err != nil {
		// Attempt to clean up the uploaded file if database insertion fails
		_ = s.storage.Delete(ctx, s.defaultBucket, image.ObjectKey)
		return nil, err
	}

	return image, nil
}

func GetImageMetadata(data []byte) (width, height int, err error) {
	config, _, err := image.DecodeConfig(strings.NewReader(string(data)))
	if err != nil {
		// Attempt to decode as base64 if initial decoding fails
		decodedData, decodeErr := base64.StdEncoding.DecodeString(string(data))
		if decodeErr != nil {
			return 0, 0, err
		}

		config, _, err = image.DecodeConfig(strings.NewReader(string(decodedData)))
		if err != nil {
			return 0, 0, err
		}
	}

	return config.Width, config.Height, nil
}
