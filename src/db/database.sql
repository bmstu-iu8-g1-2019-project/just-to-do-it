CREATE DATABASE postgres;
CREATE TABLE usertab (
    id int SERIAL PRIMARY KEY,
    email varchar(64),
    login varchar(32) UNIQUE,
    fullname varchar(128),
    password varchar(128),
    acc_verified boolean
)

CREATE TABLE auth_confirmation (
    login varchar(32),
    hash varchar(128) UNIQUE,
    deadline timestamp
)