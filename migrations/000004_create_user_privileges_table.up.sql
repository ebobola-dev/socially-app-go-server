CREATE TABLE user_privileges (
    user_id CHAR(36) NOT NULL,
    privilege_id CHAR(36) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (user_id, privilege_id),
    CONSTRAINT fk_user_privileges_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_privileges_privilege FOREIGN KEY (privilege_id) REFERENCES privileges(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_privileges_user_id ON user_privileges(user_id);
CREATE INDEX idx_user_privileges_privilege_id ON user_privileges(privilege_id);
