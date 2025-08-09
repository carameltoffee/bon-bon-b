-- +goose Up
-- +goose StatementBegin

CREATE TYPE schedule_type AS ENUM ('slot', 'deadline', 'asap');

CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type schedule_type NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS schedule_times (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    
    day_of_week VARCHAR(10) CHECK (
        day_of_week IS NULL OR day_of_week IN (
            'monday', 'tuesday', 'wednesday', 'thursday',
            'friday', 'saturday', 'sunday'
        )
    ),
    date DATE,
    slot TIME NOT NULL,
    CHECK (
        (day_of_week IS NOT NULL AND date IS NULL) OR
        (day_of_week IS NULL AND date IS NOT NULL)
    ),
    UNIQUE (schedule_id, day_of_week, date, slot)
);

CREATE TABLE IF NOT EXISTS new_schedule_slots (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    day_of_week VARCHAR(10) NOT NULL CHECK (day_of_week IN (
        'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'
    )),
    slot TIME NOT NULL,
    UNIQUE (schedule_id, day_of_week, slot)
);

CREATE TABLE IF NOT EXISTS date_slots (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    slot TIME NOT NULL,
    UNIQUE (schedule_id, date, slot)
);

CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,    
    master_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,  
    name VARCHAR(255) NOT NULL,             
    description TEXT,                       
    status VARCHAR(50) DEFAULT 'pending',  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS asap_schedules (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS deadline_schedules (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
    schedule_id INTEGER NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    deadline DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS days_off_dates (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    PRIMARY KEY (user_id, date)
);

INSERT INTO schedules (user_id, type)
SELECT DISTINCT user_id, 'slot'::schedule_type FROM schedule_slots;

INSERT INTO new_schedule_slots (schedule_id, day_of_week, slot)
SELECT s.id, ss.day_of_week, ss.slot
FROM schedule_slots ss
JOIN schedules s ON s.user_id = ss.user_id AND s.type = 'slot';

DROP TABLE IF EXISTS schedule_slots;
ALTER TABLE new_schedule_slots RENAME TO schedule_slots;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS old_schedule_slots (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    day_of_week VARCHAR(10) NOT NULL CHECK (day_of_week IN (
        'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'
    )),
    slot TIME NOT NULL
);

INSERT INTO old_schedule_slots (user_id, day_of_week, slot)
SELECT s.user_id, ss.day_of_week, ss.slot
FROM schedule_slots ss
JOIN schedules s ON s.id = ss.schedule_id AND s.type = 'slot';

DROP TABLE IF EXISTS schedule_slots;
ALTER TABLE old_schedule_slots RENAME TO schedule_slots;

DROP TABLE IF EXISTS deadline_schedules;
DROP TABLE IF EXISTS asap_schedules;
DROP TABLE IF EXISTS date_slots;
DROP TABLE IF EXISTS schedule_times;
DROP TABLE IF EXISTS days_off_dates;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS services;
DROP TYPE IF EXISTS schedule_type;

-- +goose StatementEnd