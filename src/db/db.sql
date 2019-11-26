CREATE DATABASE postgres;

CREATE TABLE user_table (
    id SERIAL PRIMARY KEY UNIQUE,
    email varchar(64) UNIQUE,
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
    deadline int,
    priority int,
    creation_datetime int,
    group_id int
);

CREATE TABLE label_table (
    id SERIAL PRIMARY KEY UNIQUE,
    task_id int,
    title varchar(32),
    color varchar(32)
);

CREATE TABLE group_table (
    id SERIAL PRIMARY KEY UNIQUE,
    title varchar(64),
    description varchar(128)
);

CREATE TABLE checklist_table (
    id SERIAL PRIMARY KEY UNIQUE,
    task_id int,
    name varchar(64)
);

CREATE TABLE checklistItem_table (
    id SERIAL PRIMARY KEY UNIQUE,
    checklist_id int,
    name varchar(64),
    state varchar(32)
);