CREATE DATABASE postgres;
CREATE TABLE user_table (
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

CREATE TABLE task_table (
    id SERIAL PRIMARY KEY UNIQUE,
    creator_id int,
    assignee_id int,
    title varchar(128),
    description varchar(128),
    state varchar(32),
    deadline TIMESTAMP WITH TIME ZONE,
    priority int,
    creation_datetime TIMESTAMP WITH TIME ZONE,
    group_id int
);

CREATE TABLE group_table (
    id SERIAL PRIMARY KEY UNIQUE,
    title varchar(32),
    description varchar(128)
);