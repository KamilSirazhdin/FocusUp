CREATE ROLE ksimouse WITH LOGIN PASSWORD 'kamilfocusup' CREATEDB;
CREATE DATABASE focusup_db OWNER ksimouse;

\c focusup_db postgres;

-- Явно установить ksimouse владельцем схемы public
ALTER SCHEMA public OWNER TO ksimouse;
GRANT ALL ON SCHEMA public TO ksimouse;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ksimouse;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ksimouse;