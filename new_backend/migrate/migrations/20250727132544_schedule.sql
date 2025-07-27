-- +goose Up
-- +goose StatementBegin
CREATE TABLE schedule_slots (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_schedule_slots_user_id ON schedule_slots(user_id);
CREATE INDEX idx_schedule_slots_time ON schedule_slots(time);
CREATE INDEX idx_schedule_slots_user_id_time ON schedule_slots(user_id, time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_schedule_slots_user_id_time;
DROP INDEX IF EXISTS idx_schedule_slots_time;
DROP INDEX IF EXISTS idx_schedule_slots_user_id;

DROP TABLE IF EXISTS schedule_slots;
-- +goose StatementEnd
