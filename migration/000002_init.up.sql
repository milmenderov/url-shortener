CREATE TABLE IF NOT EXISTS users(
    id          SERIAL not null unique,
    username    varchar(255) not null unique,
    password    varchar(255) not null);

ALTER TABLE url ADD COLUMN user_id INTEGER references users(id);