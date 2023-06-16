CREATE TABLE news_feed
(
    id    bigserial PRIMARY KEY,
    author   varchar(255) NOT NULL,
    title    varchar(255) NOT NULL,
    description TEXT NOT NULL,
    url varchar(255) NOT NULL,
    published_at TIMESTAMP NOT NULL
);