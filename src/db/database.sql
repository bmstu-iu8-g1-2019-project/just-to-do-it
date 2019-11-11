CREATE DATABASE JustToDoIt;
CREATE TABLE task_table (
    id SERIAL PRIMARY KEY UNIQUE,
    assignee_id int,
    title varchar(128),
    description varchar(128),
    state varchar(32),
    deadline TIMESTAMP WITH TIME ZONE,
    priority int,
    creation_datetime TIMESTAMP WITH TIME ZONE,
    group_id int
)