-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    points INT DEFAULT 0,
    streak INT DEFAULT 0,
    last_active_date TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS tasks (
                                     id SERIAL PRIMARY KEY,
                                     question TEXT NOT NULL,
                                     answer TEXT NOT NULL,
                                     points INT DEFAULT 10,
                                     created_by_id INT,
                                     created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS user_task_logs (
                                              id SERIAL PRIMARY KEY,
                                              user_id INT NOT NULL,
                                              task_id INT NOT NULL,
                                              answered_at TIMESTAMP NOT NULL,
                                              correct BOOLEAN DEFAULT false,
                                              created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
    );

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS user_task_logs;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;