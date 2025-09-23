package postgres

import (
	"context"
	"database/sql"

	"github.com/viniferr33/img-processor/internal/image"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) GetByID(ctx context.Context, id string) (*image.Image, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, original_name, size, mime_type, width, height, parent_id, version, created_at, updated_at, bucket_name, object_key
		FROM images
		WHERE id = $1
	`, id)

	var img image.Image
	err := row.Scan(
		&img.ID,
		&img.OriginalName,
		&img.Size,
		&img.MimeType,
		&img.Width,
		&img.Height,
		&img.ParentID,
		&img.Version,
		&img.CreatedAt,
		&img.UpdatedAt,
		&img.BucketName,
		&img.ObjectKey,
	)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *ImageRepository) Create(ctx context.Context, img *image.Image) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO images (id, original_name, size, mime_type, width, height, parent_id, version, created_at, updated_at, bucket_name, object_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW()), $9, $10)
	`, img.ID, img.OriginalName, img.Size, img.MimeType, img.Width, img.Height, img.ParentID, img.Version, img.BucketName, img.ObjectKey)
	return err
}

func (r *ImageRepository) Update(ctx context.Context, img *image.Image) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE images
		SET original_name = $1,
			size = $2,
			mime_type = $3,
			width = $4,
			height = $5,
			parent_id = $6,
			version = $7,
			updated_at = EXTRACT(EPOCH FROM NOW()),
			bucket_name = $8,
			object_key = $9
		WHERE id = $10
	`, img.OriginalName, img.Size, img.MimeType, img.Width, img.Height, img.ParentID, img.Version, img.BucketName, img.ObjectKey, img.ID)
	return err
}

func (r *ImageRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM images
		WHERE id = $1
	`, id)
	return err
}
