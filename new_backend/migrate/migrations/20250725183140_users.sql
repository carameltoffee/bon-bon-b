-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    bio TEXT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    specialization VARCHAR(100) DEFAULT '-',
    role CHAR(5) NOT NULL DEFAULT 'user',                 
    is_active BOOLEAN NOT NULL DEFAULT TRUE,          
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,     
    last_login_at TIMESTAMP,                        
    last_ip TEXT,                                     
    deleted_at TIMESTAMP,                          
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
