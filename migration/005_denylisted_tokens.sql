CREATE TABLE IF NOT EXISTS denylist (
    token       TEXT PRIMARY KEY,  
    expires_at  TIMESTAMP NOT NULL
);