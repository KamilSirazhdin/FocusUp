-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Добавление администратора в систему
INSERT INTO users (
    username, email, password_hash, role, points, streak, last_active_date, created_at, updated_at
) VALUES (
             'admin', 'admin@focusup.com',
             '$2b$14$mj3ZzC7mtGcuVs7KaiglResCANeyMyJk.xNzwy4AkLWS9pbneS4Ua', -- пароль: admin123
             'admin', 0, 0, NOW(), NOW(), NOW()
         ) ON CONFLICT DO NOTHING;

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_points_desc ON users(points DESC);
CREATE INDEX IF NOT EXISTS idx_user_task_logs_user_id ON user_task_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_task_logs_task_id ON user_task_logs(task_id);
CREATE INDEX IF NOT EXISTS idx_user_task_logs_answered_at ON user_task_logs(answered_at DESC);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

-- Удаление индексов
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_points_desc;
DROP INDEX IF EXISTS idx_user_task_logs_user_id;
DROP INDEX IF EXISTS idx_user_task_logs_task_id;
DROP INDEX IF EXISTS idx_user_task_logs_answered_at;

-- Удаление администратора
DELETE FROM users WHERE email = 'admin@focusup.com';
