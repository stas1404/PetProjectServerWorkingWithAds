CREATE TABLE IF NOT EXISTS ads (
    id serial not null primary key,
    title varchar(100) not null,
    text varchar(400) not null,
    author_id serial not null,
    published bool default false,
    created date,
    last_modified date
);