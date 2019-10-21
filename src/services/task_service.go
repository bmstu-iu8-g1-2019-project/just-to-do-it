package services

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

const (
	g_query = "(assignee_id, title, description, state," +
		"deadline, priority, creation_datetime, group_id) values" +
		"($1, $2, $3, $4, $5, $6, $7, $8, $9)"
)


type DatastoreTask interface  {
	GetTaskTId(id int64) (*models.Task)
	GetTasksAId(id int64) []models.Task
	GetTasksGId(id int64) []models.Task
	UpdateTask(task models.Task) error
	CreateTask(task models.Task) error
}


func (db *DB) GetTaskTId(id int64) *models.Task {
	row := db.QueryRow("SELECT * FROM task where id = $1")

	defer db.Close()

	task := &models.Task{}
	err := row.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
		&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
	if err != nil {
		return nil
	}
	return task
}

func (db *DB) GetTasksAId(id int64) []models.Task {
	rows, err := db.Query("SELECT * FROM task where assignee_id = $1", id)

	defer db.Close()

	if err != nil {
		return nil
	}

	tasks := make([]models.Task, 0)

	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
			&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
		if err != nil {
			return nil
		}

		tasks = append(tasks, *task)
	}
	return tasks
}

func (db *DB) GetTasksGId(id int64) []models.Task {
	rows, err := db.Query("SELECT * FROM task where group_id = $1", id)

	defer db.Close()

	if err != nil {
		return nil
	}

	tasks := make([]models.Task, 0)

	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
			&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
		if err != nil {
			return nil
		}

		tasks = append(tasks, *task)
	}
	return tasks
}

func (db *DB) UpdateTask(task models.Task) error {
	_, err := db.Exec("UPDATE task SET " + g_query, task.Id, task.AssigneeId, task.Title,
		task.Description,
			task.State, task.Deadline, task.Priority, task.CreationDatetime,
				task.GroupId)

	defer db.Close()

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateTask(task models.Task) error {
	_, err := db.Exec("INSERT INTO task " + g_query, task.Id, task.AssigneeId, task.Title,
		task.Description,
		task.State, task.Deadline, task.Priority, task.CreationDatetime,
		task.GroupId)

	defer db.Close()

	if err != nil {
		return err
	}
	return nil
}


