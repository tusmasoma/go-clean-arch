CREATE DATABASE IF NOT EXISTS `goCleanArcTestDB` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `goCleanArcTestDB`;

DROP TABLE IF EXISTS Tasks CASCADE;
DROP TABLE IF EXISTS Users CASCADE;

-- Tasks Table
CREATE TABLE Tasks (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duedate TIMESTAMP,
    priority INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users Table
CREATE TABLE Users (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);