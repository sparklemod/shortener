-- +goose Up
-- +goose StatementBegin
ALTER TABLE links
    RENAME COLUMN redirect_url TO shorten_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links
    RENAME COLUMN redirect_url TO shorten_url;
-- +goose StatementEnd
