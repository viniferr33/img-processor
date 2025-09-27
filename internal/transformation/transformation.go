package transformation

import "time"

type Pipeline struct {
	ID             string
	ImageID        string
	CreatedAt      time.Time
	TransformSteps []TransformStep
}

type TransformStep struct {
	ID         string
	PipelineID string
	Order      int
	Type       TransformationType
	Params     map[string]string // JSON-encoded params for flexibility
}

type TransformationType string

const (
	Resize TransformationType = "resize"
	Crop   TransformationType = "crop"
	Rotate TransformationType = "rotate"
	Filter TransformationType = "filter"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type PipelineRepository interface {
	Create(pipeline *Pipeline) error
	GetByID(id int64) (*Pipeline, error)
	Update(pipeline *Pipeline) error
	Delete(id int64) error
}
