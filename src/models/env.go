package models

import (
	"database/sql"
)

type Datastore interface {
	GetTaskByAssigneeId(id int) []Task
	GetTaskByGroupId(id int) []Task
	GetTaskByTaskId(id int) *Task
}

type Environment struct {
	db Datastore
}

type DB struct {
	*sql.DB
}
