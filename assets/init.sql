CREATE TABLE users (
    id UUID PRIMARY KEY,
    email text UNIQUE,
    login text UNIQUE,
    fullname text,
    password text,
    acc_verified boolean
);

CREATE TABLE auth_confirmation (
    login text UNIQUE,
    hash text UNIQUE,
    deadline bigint
);

CREATE TABLE tasks (
    uuid UUID PRIMARY KEY,
    creator_id UUID,
    assignee_id UUID,
    group_id UUID,
    title text,
    description text,
    state text,
    deadline integer,
    duration integer,
    priority integer,
    creation_datetime bigint
);

CREATE TABLE groups (
    uuid UUID PRIMARY KEY,
    creator_id UUID,
    title text,
    description text
);

CREATE TABLE labels (
    uuid UUID PRIMARY KEY,
    task_id UUID,
    title text,
    color text
);

CREATE TABLE tracks (
    uuid UUID PRIMARY KEY,
    group_id UUID,
    title text,
    description text
);

CREATE TABLE track_item (
    uuid UUID UNIQUE,
    previous_id UUID,
    track_id UUID
);

CREATE TABLE checklists (
    uuid UUID PRIMARY KEY,
    task_id UUID,
    name text
);

CREATE TABLE checklist_items (
    uuid UUID PRIMARY KEY,
    checklist_id UUID,
    name text,
    state text
);

CREATE TABLE scopes (
    uuid UUID PRIMARY KEY,
    creator_id UUID,
    group_id UUID,
    begin_interval bigint,
    end_interval bigint
);

CREATE TABLE timetables (
    scope_id UUID,
    task_id UUID
);
