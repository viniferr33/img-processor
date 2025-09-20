package image

import "context"

type ImageService struct {
	repo    ImageRepository
	storage ObjectStorage
}

func NewImageService(repo ImageRepository, storage ObjectStorage) *ImageService {
	return &ImageService{
		repo:    repo,
		storage: storage,
	}
}

func (s *ImageService) UploadImage(ctx context.Context, data []byte, filename string) (*Image, error) {
	return nil, nil
}
