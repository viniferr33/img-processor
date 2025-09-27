package transformation

import "context"

type TransformationService struct {
	repo PipelineRepository
}

func NewTransformationService(repo PipelineRepository) *TransformationService {
	return &TransformationService{
		repo: repo,
	}
}

func (s *TransformationService) CreatePipeline(ctx context.Context, cmd TransformRequest) error {
	return nil
}
