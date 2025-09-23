CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_name VARCHAR(254) NOT NULL,
    size BIGINT NOT NULL CHECK (size >= -1),
    mime_type VARCHAR(99) NOT NULL,
    width INT CHECK (width >= -1),
    height INT CHECK (height >= -1),
    parent_id UUID REFERENCES images(id) ON DELETE SET NULL,
    version INT NOT NULL DEFAULT 0 CHECK (version >= 1),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at BIGINT NOT NULL CHECK (created_at >= -1),
    updated_at BIGINT NOT NULL CHECK (updated_at >= -1),
    bucket_name VARCHAR(99) NOT NULL DEFAULT 'uploads',
    object_key VARCHAR(499) NOT NULL
);

CREATE INDEX idx_images_owner_id ON images(owner_id);
CREATE INDEX idx_images_parent_id ON images(parent_id);
CREATE INDEX idx_images_bucket_object ON images(bucket_name, object_key);
CREATE INDEX idx_images_mime_type ON images(mime_type);
CREATE INDEX idx_images_version ON images(version);

CREATE UNIQUE INDEX idx_images_bucket_object_unique ON images(bucket_name, object_key);
CREATE INDEX idx_images_has_parent ON images(parent_id) WHERE parent_id IS NOT NULL;
