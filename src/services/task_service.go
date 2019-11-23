package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

type DatastoreTask interface {
	GetTasks([]int, string, int) ([]models.Task, error)
	//GetTask(int, int, string, int) (models.Task, error)
	//UpdateTask(models.Task, int) error
	CreateTask(models.Task) (models.Task, error)
	//DeleteTask(int) error
}

// input we get an array of values from url
// returns an array of objects of type Task
func(db* DB) GetTasks(idSlice []int, title string, userId int) (tasks []models.Task ,err error) {
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
	if userId != 0 {
		queryMap["creator_id"]= userId
	}
	query := "SELECT id, creator_id, assignee_id, title, description, state, deadline, priority, creation_datetime, group_id FROM task_table WHERE "

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
		err = rows.Scan(&task.Id, &task.CreatorId, &task.AssigneeId, &task.Title, &task.Description,
			&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, *task)
	}
	return tasks, nil
}

//create task
func (db *DB) CreateTask(task models.Task) (models.Task, error) {
	_, err := db.Exec("INSERT INTO task_table (creator_id, assignee_id, title, description, state, deadline, priority," +
		"creation_datetime, group_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		task.CreatorId, task.AssigneeId,
		task.Title, task.Description,
		task.State, task.Deadline,
		task.Priority, time.Now(),
		task.GroupId)

	if err != nil {
		return task, err
	}
	return task,nil
}
