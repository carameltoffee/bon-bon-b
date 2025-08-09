-- +goose Up
-- +goose StatementBegin
CREATE TABLE schedule_patterns (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    schedule_type VARCHAR(20) NOT NULL, 
    days_ahead INTEGER NOT NULL DEFAULT 30,
    last_generated TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE schedule_pattern_time_ranges (
    pattern_id BIGINT NOT NULL REFERENCES schedule_patterns(id) ON DELETE CASCADE,
    weekday SMALLINT NOT NULL CHECK (weekday BETWEEN 0 AND 6) ,
    time TIME NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schedule_pattern_time_ranges;
DROP TABLE IF EXISTS schedule_patterns;
-- +goose StatementEnd
