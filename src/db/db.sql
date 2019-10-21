CREATE DATABASE postgres;
CREATE TABLE usertab (
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
