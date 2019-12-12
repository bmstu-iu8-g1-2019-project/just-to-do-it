drop table if exists user_table cascade;
drop table if exists auth_confirmation cascade;
drop table if exists group_table cascade;
drop table if exists label_table cascade;

CREATE TABLE if not exists user_table (
    id SERIAL PRIMARY KEY UNIQUE,
    email varchar(64) UNIQUE,
    login varchar(32) UNIQUE,
    fullname varchar(128),
    password varchar(128),
    acc_verified boolean
);

CREATE TABLE if not exists auth_confirmation (
    login varchar(32) UNIQUE,
    hash varchar(128) UNIQUE,
    deadline TIMESTAMP WITH TIME ZONE
);

CREATE TABLE if not exists task_table (
    id SERIAL PRIMARY KEY UNIQUE,
    creator_id int,
    assignee_id int,
    title varchar(128),
    description varchar(128),
    state varchar(32),
    deadline int,
    duration int,
    priority int,
    creation_datetime int,
    group_id int
);

CREATE TABLE if not exists group_table (
    id SERIAL PRIMARY KEY UNIQUE,
    creator_id int,
    title varchar(32),
    description varchar(128)
);

CREATE TABLE if not exists label_table (
    id SERIAL PRIMARY KEY UNIQUE,
    task_id int,
    title varchar(32),
    color varchar(32)
);

CREATE TABLE if not exists track_table (
    id SERIAL PRIMARY KEY UNIQUE,
    title varchar(32),
    description varchar(128),
    group_id int
);

CREATE TABLE if not exists track_task_previous (
    task_id int UNIQUE,
    previous_id int,
    track_id int
);

CREATE TABLE if not exists checklist_table (
    id SERIAL PRIMARY KEY UNIQUE,
    task_id int,
    name varchar(128)
);

CREATE TABLE if not exists checklistItem_table (
    id SERIAL PRIMARY KEY UNIQUE,
    checklist_id int,
    name varchar(128),
    state varchar(32)
);

CREATE TABLE if not exists scope (
    id SERIAL PRIMARY KEY UNIQUE,
    creator_id int,
    group_id int,
    begin_interval int,
    end_interval int
);

CREATE TABLE if not exists timetable (
    scope_id int,
    task_id int
);
