CREATE TYPE expense_category AS ENUM (
    'purchase',         -- покупка объекта
    'rough_work',       -- черновые работы
    'finish_work',      -- чистовые работы
    'materials_rough',  -- материалы (черновые)
    'materials_finish', -- материалы (чистовые)
    'plumbing',         -- сантехника
    'electrical',       -- электрика
    'furniture',        -- мебель и техника
    'electronic'
    'commission',       -- комиссия агента
    'taxes',            -- налоги
    'other'             -- прочее
);

CREATE TABLE IF NOT EXISTS expenses (
    id          UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID             NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id     UUID             NOT NULL REFERENCES users(id),
    category    expense_category NOT NULL,
    amount      NUMERIC(15, 2)   NOT NULL CHECK (amount >= 0),
    description TEXT,
    paid_at     DATE             NOT NULL DEFAULT CURRENT_DATE,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expenses_project_id ON expenses(project_id);
CREATE INDEX idx_expenses_category    ON expenses(category);
CREATE INDEX idx_expenses_paid_at     ON expenses(paid_at);