package services

import (
	"strconv"
	"strings"
	"fmt"
	"time"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

type DatastoreTask interface  {
	GetTasks([]int, string) ([]models.Task, error)
	GetTask(int, int, string, int) (models.Task, error)
	UpdateTask(models.Task, int) error
	CreateTask(models.Task) error
	DeleteTask(int) error
}

func(db* DB) GetTask(id int, assigneeId int, title string, groupId int) (task models.Task ,err error) {
	queryMap := make(map[string]interface{})
	if id != 0 {
		queryMap["id"] = id
	}
	if assigneeId != 0 {
		queryMap["assignee_id"] = assigneeId
	}
	if title != "" {
		queryMap["title"] = title
	}
	if groupId != 0 {
		queryMap["group_id"] = groupId
	}
	query := "SELECT id, assignee_id, title, description, state, deadline, priority, creation_datetime, group_id FROM task_table WHERE "

	var values []interface{}
	var where []string
	i := 1
	for k, v := range queryMap {
		values = append(values, v)
		where = append(where, fmt.Sprintf("%s = $%s", k, strconv.Itoa(i)))
		i++
	}
	//fmt.Println(query + strings.Join(where, " AND "))

	row := db.QueryRow(query + strings.Join(where, " AND "), values...)
	if err != nil {
		return task, err
	}

	err = row.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
		&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
		&task.GroupId)
	if err != nil {
		return task, err
	}

	return task, nil
}

func(db* DB) GetTasks(idSlice []int, title string) (tasks []models.Task ,err error) {
	queryMap := make(map[string]interface{})
	if idSlice[0] != 0 {
		queryMap["id"] = idSlice[0]
	}
	if idSlice[1] != 0 {
		queryMap["assignee_id"] = idSlice[1]
	}
	if idSlice[2] != 0 {
		queryMap["group_id"] = idSlice[2]
	}
	if title != "" {
		queryMap["title"] = title
	}
	query := "SELECT id, assignee_id, title, description, state, deadline, priority, creation_datetime, group_id FROM task_table WHERE "

	var values []interface{}
	var where []string
	i := 1
	for k, v := range queryMap {
		values = append(values, v)
		where = append(where, fmt.Sprintf("%s = $%s", k, strconv.Itoa(i)))
		i++
	}
	//fmt.Println(query + strings.Join(where, " AND "))

	rows, err := db.Query(query + strings.Join(where, " AND "), values...)
	if err != nil {
		return tasks, err
	}

	tasks = make([]models.Task, 0)

	for rows.Next() {
		task := &models.Task{}
		err = rows.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
			            &task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			            &task.GroupId)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (db *DB) UpdateTask(task models.Task, Id int) error {
	_, err := db.Exec("UPDATE task_table SET assignee_id = $1, title = $2, description = $3, state = $4, deadline = $5," +
		" priority = $6, creation_datetime = $7, group_id = $8 where id = $9",
		task.AssigneeId, task.Title, task.Description, task.State,
		task.Deadline, task.Priority, task.CreationDatetime,
		task.GroupId, Id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetTaskById (id int) (task models.Task, err error) {
	row := db.QueryRow("SELECT * FROM task_table WHERE id = $1", id)
	err = row.Scan(&task.Id, &task.AssigneeId, &task.Title, &task.Description,
		&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime, &task.GroupId)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (db *DB) CreateTask(task models.Task) error {
	_, err := db.Exec("INSERT INTO task_table (assignee_id, title, description, state, deadline, priority," +
		"creation_datetime, group_id) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		task.AssigneeId, task.Title,
		task.Description,
		task.State, time.Time{}, task.Priority, time.Now(),
		task.GroupId)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteTask(id int) (err error) {
	_, err = db.GetTaskById(id)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM task_table WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
