CREATE TABLE IF NOT EXISTS users (
    id serial not null primary key,
    nickname VARCHAR(25) not null,
    email VARCHAR(25) not null,
    password VARCHAR(25) not null
);