-- +goose Up
ALTER TABLE user_task_logs ADD COLUMN deleted_at TIMESTAMP;

-- +goose Down
ALTER TABLE user_task_logs DROP COLUMN deleted_at;