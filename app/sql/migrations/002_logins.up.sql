CREATE TABLE logins (
    id VARCHAR PRIMARY KEY,
    password_hash VARCHAR NOT NULL,
    last_login TIMESTAMP,
    failed_attempts INT DEFAULT 0,
    is_locked BOOLEAN DEFAULT FALSE,
    user_id VARCHAR NOT NULL,
    CONSTRAINT fk_logins_user 
        FOREIGN KEY(user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);