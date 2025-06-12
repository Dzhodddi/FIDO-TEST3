CREATE TABLE IF NOT EXISTS quotes (
    id bigserial PRIMARY KEY,
    author varchar(200) NOT NULL,
    text TEXT NOT NULL
);