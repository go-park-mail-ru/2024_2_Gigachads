-- Создание таблицы пользователей (profile)
CREATE TABLE IF NOT EXISTS profile (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email TEXT NOT NULL UNIQUE CHECK (LENGTH(email) <= 50),
    username TEXT NOT NULL CHECK (LENGTH(username) >= 5 AND LENGTH(username) <= 50),
    password TEXT NOT NULL CHECK (LENGTH(password) >= 5 AND LENGTH(password) <= 50),
    avatar_url TEXT DEFAULT 'default.png'
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
    CONSTRAINT fk_message FOREIGN KEY (message_id) REFERENCES message(id),
    CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES folder(id) ON DELETE CASCADE
);

-- Создание таблицы вложений
CREATE TABLE IF NOT EXISTS attachment (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    message_id INTEGER,
    url TEXT NOT NULL,
    CONSTRAINT fk_message FOREIGN KEY  (message_id) REFERENCES message(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS question (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    action TEXT IN ('Send', 'Answer', 'Forward', 'UploadAvatar', 'SignUp', 'Main', 'Delete')
    type TEXT IN ('Star', 'Number')
    description TEXT
)

CREATE TABLE IF NOT EXISTS answer (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    action TEXT IN ('Send', 'Answer', 'Forward', 'UploadAvatar', 'SignUp', 'Main', 'Delete')
    user_email TEXT, 
    CONSTRAINT fk_user FOREIGN KEY (user_email) REFERENCES profile(email) ON DELETE NULL,
    value INT
)

INSERT INTO question (action, type, dexcription) VALUES ('Send', 'Как бы вы оценили качество работы почтового сервиса?', );
INSERT INTO question (action, type, dexcription) VALUES ('Forward', 'Number', 'Оцените по шкале от 1 до 10, насколько просто вам пересылать сообщения (где 1 - очень сложно, 10- очень легко).');
INSERT INTO question (action, type, dexcription) VALUES ('Answer', 'Number', 'Оцените по шкале от 1 до 10, насколько просто вам отвечать на сообщения (где 1 - очень сложно, 10- очень легко).');
INSERT INTO question (action, type, dexcription) VALUES ('UploadAvatar', 'Насколько просто и удобно вам менять аватарку в почтовом сервисе?');
INSERT INTO question (action, type, dexcription) VALUES ('SignUp', 'Насколько легким и понятным вы считаете процесс регистрации на сайте Gigamail?');
INSERT INTO question (action, type, dexcription) VALUES ('Main', 'Как вы оцениваете удобство навигации в почтовом сервисе?');
INSERT INTO question (action, type, dexcription) VALUES ('Delete', 'Number', 'Оцените по шкале от 1 до 10, насколько легко вам удалять письма в почтовом сервисе (где 1 - очень сложно, 10- очень легко).');

    
    
    
    
    
    
    