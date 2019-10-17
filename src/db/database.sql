CREATE DATABASE JustToDoIt;
CREATE TABLE task (
    id int SERIAL PRIMARY KEY,
    assignee_id int,
    title text(128),
    description string,
    state varchar(32),
    deadline time,
    priority int,
    creation_datetime timestamp,
    group_id int
)