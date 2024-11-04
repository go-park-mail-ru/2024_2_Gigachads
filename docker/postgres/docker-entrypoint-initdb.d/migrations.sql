-- CREATE DATABASE IF NOT EXISTS test;

-- Создание таблицы пользователей (profile)
CREATE TABLE IF NOT EXISTS profile (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email TEXT NOT NULL UNIQUE CHECK (LENGTH(email) <= 50),
    username TEXT NOT NULL CHECK (LENGTH(username) >= 5 AND LENGTH(username) <= 50),
    password TEXT NOT NULL CHECK (LENGTH(password) >= 5 AND LENGTH(password) <= 50),
    avatar_url TEXT DEFAULT NULL
);

-- Создание таблицы писем (message)
CREATE TABLE IF NOT EXISTS message (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL CHECK (LENGTH(title) <= 100),
    description TEXT DEFAULT NULL
);

-- Создание таблицы папок (folder)
CREATE TABLE IF NOT EXISTS folder (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER,
    name TEXT NOT NULL CHECK (LENGTH(name) <= 50),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES profile(id) ON DELETE CASCADE
);

-- Создание таблицы писем пользователей (email_transaction)
CREATE TABLE IF NOT EXISTS email_transaction (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    parent_transaction_id INTEGER DEFAULT NULL,
    sender_email TEXT,
    recipient_email TEXT,
    message_id INTEGER,
    sending_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    isRead BOOLEAN NOT NULL DEFAULT FALSE,
    folder_id INTEGER,
    CONSTRAINT fk_sender FOREIGN KEY (sender_email) REFERENCES profile(email) ON DELETE SET NULL,
    CONSTRAINT fk_recipient FOREIGN KEY (recipient_email) REFERENCES profile(email) ON DELETE SET NULL,
    CONSTRAINT fk_message FOREIGN KEY (message_id)  REFERENCES message(id),
    CONSTRAINT fk_folder FOREIGN KEY (folder_id)  REFERENCES folder(id) ON DELETE CASCADE
);

-- Создание таблицы вложений
CREATE TABLE IF NOT EXISTS attachment (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    message_id INTEGER,
    url TEXT NOT NULL,
    CONSTRAINT fk_message FOREIGN KEY  (message_id) REFERENCES message(id) ON DELETE CASCADE
);