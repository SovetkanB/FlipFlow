CREATE TYPE file_entity_type AS ENUM ('project', 'expense');
CREATE TYPE file_type AS ENUM ('photo', 'drawing', 'document', 'other');

CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type file_entity_type NOT NULL,
    entity_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    file_type file_type NOT NULL DEFAULT 'photo',
    original_name VARCHAR(255) NOT NULL,
    storage_key TEXT NOT NULL UNIQUE,
    bucket VARCHAR(100) NOT NULL,
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_files_entity ON files(entity_type, entity_id);
CREATE INDEX idx_files_user   ON files(user_id);
