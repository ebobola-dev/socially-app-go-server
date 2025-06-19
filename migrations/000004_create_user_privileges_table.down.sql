-- Сначала внешние ключи
ALTER TABLE user_privileges DROP FOREIGN KEY fk_user_privileges_user;
ALTER TABLE user_privileges DROP FOREIGN KEY fk_user_privileges_privilege;

-- Потом вручную созданные индексы (уже безопасно)
DROP INDEX idx_user_privileges_user_id ON user_privileges;
DROP INDEX idx_user_privileges_privilege_id ON user_privileges;

-- И только потом таблицу
DROP TABLE IF EXISTS user_privileges;