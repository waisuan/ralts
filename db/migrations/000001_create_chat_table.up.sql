CREATE TABLE chat
(
    chat_id    bigserial PRIMARY KEY,
    username   varchar(255) NOT NULL,
    message    text      NOT NULL,
    created_at TIMESTAMP NOT NULL
);