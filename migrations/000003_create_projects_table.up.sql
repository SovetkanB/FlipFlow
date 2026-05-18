CREATE TYPE project_status AS ENUM (
    'search',
    'purchased',
    'renovation',
    'for_sale',
    'sold'
);

CREATE TABLE IF NOT EXISTS projects (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    address         TEXT,
    city            VARCHAR(100),
    area_sqm        NUMERIC(8, 2),
    rooms           SMALLINT,
    purchase_price  NUMERIC(15, 2),
    target_price    NUMERIC(15, 2),
    sold_price      NUMERIC(15, 2),
    status project_status NOT NULL DEFAULT 'search',
    description     TEXT,
    purchased_at    DATE,
    sold_at         DATE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_projects_user_id  ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_projects_status   ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_city     ON projects(city);