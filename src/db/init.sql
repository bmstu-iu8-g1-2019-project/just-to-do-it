

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
    title varchar(32),
    description varchar(128)
);

CREATE TABLE if not exists label_table (
    id SERIAL PRIMARY KEY UNIQUE,
    task_id int,
    title varchar(32),
    color varchar(32)
);
