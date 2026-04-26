CREATE TABLE logins (
    id UUID PRIMARY KEY,
    password_hash VARCHAR NOT NULL,
    last_login TIMESTAMP DEFAULT NULL,
    failed_attempts INT DEFAULT 0,
    is_locked BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    CONSTRAINT fk_logins_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);