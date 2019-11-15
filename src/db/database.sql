CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY,
    assignee_id int,
    title TEXT DEFAULT '',
    description TEXT DEFAULT '',
    state varchar(32),
    deadline time,
    priority int,
    creation_datetime timestamp,
    group_id int
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY UNIQUE,
    email varchar(64),
    login varchar(32) UNIQUE,
    fullname varchar(128),
    password varchar(128),
    acc_verified boolean
);

CREATE TABLE auth_confirmation (
    login varchar(32) UNIQUE,
    hash varchar(128) UNIQUE,
    deadline TIMESTAMP WITH TIME ZONE
);