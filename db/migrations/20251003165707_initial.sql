-- +goose Up
CREATE TABLE links (
                       id           SERIAL PRIMARY KEY,
                       original_url TEXT NOT NULL,
                       redirect_url TEXT UNIQUE NOT NULL,
                       is_active    BOOLEAN NOT NULL DEFAULT TRUE,
                       visits       INTEGER NOT NULL DEFAULT 0,
                       created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP table links;
