CREATE DATABASE goCleanArcTestDB;

\c goCleanArcTestDB;

DROP TABLE IF EXISTS Tasks CASCADE;

CREATE TABLE Tasks (
    id CHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duedate TIMESTAMP,
    priority INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
