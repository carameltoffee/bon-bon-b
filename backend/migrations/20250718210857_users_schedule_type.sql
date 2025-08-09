-- +goose Up
-- +goose StatementBegin

ALTER TABLE users ADD COLUMN default_schedule_type schedule_type DEFAULT 'slot';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users DROP COLUMN default_schedule_type;

-- +goose StatementEnd
