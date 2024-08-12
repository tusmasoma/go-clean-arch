CREATE DATABASE IF NOT EXISTS `goCleanArcDB` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `goCleanArcDB`;

DROP TABLE IF EXISTS Tasks CASCADE;

-- Tasks Table
CREATE TABLE Tasks (
    id CHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duedate TIMESTAMP,
    priority INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);