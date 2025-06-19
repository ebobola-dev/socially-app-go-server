CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(16) NOT NULL,
    password CHAR(60) NOT NULL,

    fullname VARCHAR(32),
    about_me VARCHAR(256),
    gender ENUM('male', 'female'),
    date_of_birth DATE NOT NULL,

    avatar_type ENUM(
        'external', 'avatar1', 'avatar2', 'avatar3', 'avatar4',
        'avatar5', 'avatar6', 'avatar7', 'avatar8', 'avatar9', 'avatar10'
    ),
    avatar_id CHAR(36) UNIQUE,

    last_seen DATETIME(3),
    deleted_at DATETIME(3),
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE UNIQUE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);
