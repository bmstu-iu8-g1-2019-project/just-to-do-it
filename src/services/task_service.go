package services

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"database/sql"
)

type DB struct {
	*sql.DB
}

func (db *DB) GetTaskByTaskId(id int64) (*models.Task) {
	query := "SELECT * FROM task where id = $1"
	row := db.QueryRow(query)

	task := &models.Task{}
	err := row.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
		&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
	if err != nil {
		return nil
	}
	return task
}

func (db *DB) GetTaskByAssigneeId(id int64) []models.Task {
	query := "SELECT * FROM task where assignee_id = $1"
	rows, err := db.Query(query, id)
	if err != nil {
		panic(err)
	}

	tasks := make([]models.Task, 0)

	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
			&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
		if err != nil {
			panic(err)
		}

		tasks = append(tasks, *task)
	}
	return tasks
}

func (db *DB) GetTaskByGroupId(id int64) []models.Task {
	query := "SELECT * FROM task where group_id = $1"
	rows, err := db.Query(query, id)
	if err != nil {
		panic(err)
	}

	tasks := make([]models.Task, 0)

	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
			&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
		if err != nil {
			panic(err)
		}

		tasks = append(tasks, *task)
	}
	return tasks
}

func (db *DB) UpdateTaskAllProperties(task models.Task) error {
	query := "UPDATE task SET (id, assignee_id, title, description, state," +
		"deadline, priority, creation_datetime, group_id) values" +
			"($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	_, err := db.Exec(query, task.Id, task.AssigneeId, task.Title,
		task.Description,
			task.State, task.Deadline, task.Priority, task.CreationDatetime,
				task.GroupId)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) PostTask(task models.Task) error {
	query := "INSERT INTO task (assignee_id, title, description, state," +
		"deadline, priority, creation_datetime, group_id)values" +
			"($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	_, err := db.Exec(query, task.Id, task.AssigneeId, task.Title,
		task.Description,
		task.State, task.Deadline, task.Priority, task.CreationDatetime,
		task.GroupId)

	if err != nil {
		return err
	}
	return nil
}


