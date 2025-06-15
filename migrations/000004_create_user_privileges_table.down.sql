DROP INDEX idx_user_privileges_user_id ON user_privileges;
DROP INDEX idx_user_privileges_privilege_id ON user_privileges;

DROP TABLE IF EXISTS user_privileges;
