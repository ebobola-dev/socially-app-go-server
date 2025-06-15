DROP INDEX idx_users_email ON users;
DROP INDEX idx_users_username ON users;
DROP INDEX idx_users_deleted_at ON users;

DROP TABLE IF EXISTS users;
