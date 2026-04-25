CREATE TABLE IF NOT EXISTS questions (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    topic       TEXT NOT NULL,
    difficulty  TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    slug        TEXT NOT NULL UNIQUE
);