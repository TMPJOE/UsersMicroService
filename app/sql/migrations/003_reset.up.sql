CREATE TABLE password_resets (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    reset_token VARCHAR NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_resets_user 
        FOREIGN KEY(user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);