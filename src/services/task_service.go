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
	GetTaskById(int) (models.Task, []models.Label, error)
	UpdateTask(models.Task, int, int) (models.Task, error)
	CreateTask(models.Task, int) (models.Task, error)
	//DeleteTask(int) error
	CreateLabel(models.Label, int) (models.Label, error)
	GetLabelsByTaskId(int) ([]models.Label, error)
	GetLabel(int) (models.Label, error)
	UpdateLabelColor(int, string) (models.Label, error)
	UpdateLabelTitle(int, string) (models.Label, error)
	DeleteLabel(int) error
	//
	//CreateChecklist(models.Checklist, int) (models.Checklist, error)
	//CreateChecklistItem(models.ChecklistItem, int) (models.ChecklistItem, error)
	//GetChecklist(int) (models.Checklist, []models.ChecklistItem, error)
	//UpdateChecklist(int, models.Checklist) (models.Checklist, error)
	//DeleteChecklist(int) error
	//GetChecklistItems(int) ([]models.ChecklistItem, error)
	//GetChecklistItem(int) (models.ChecklistItem, error)
	//UpdateChecklistItem(int, int, models.ChecklistItem) (models.ChecklistItem, error)
	//DeleteChecklistItem(int) error
}

// input we get an array of values from url
// returns an array of objects of type Task
func(db* DB) GetTasks(idSlice []int, title string, userId int) (tasks []models.Task, err error) {
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

	rows, err := db.Query(query + strings.Join(where, " AND "), values...)
	if err != nil {
		return []models.Task{}, err
	}

	tasks = make([]models.Task, 0)

	for rows.Next() {
		task := &models.Task{}
		err = rows.Scan(&task.Id, &task.CreatorId, &task.AssigneeId, &task.Title, &task.Description,
			&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
		if err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

//get task
func (db *DB) GetTaskById (id int) (task models.Task, labels []models.Label, err error) {
	row := db.QueryRow("SELECT * FROM task_table WHERE id = $1", id)
	err = row.Scan(&task.Id, &task.CreatorId, &task.AssigneeId, &task.Title, &task.Description,
		&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime, &task.GroupId)
	if err != nil {
		return models.Task{}, []models.Label{}, err
	}
	labels, err = db.GetLabelsByTaskId(id)
	if err != nil {
		return models.Task{}, []models.Label{}, err
	}
	return task, labels,nil
}

//create task
func (db *DB) CreateTask(task models.Task, userId int) (models.Task, error) {
	err := db.QueryRow("INSERT INTO task_table (creator_id, assignee_id, title, description, state, deadline," +
		"priority, creation_datetime, group_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)  RETURNING id",
		userId, task.AssigneeId, task.Title, task.Description, task.State,
		task.Deadline, task.Priority, time.Now().Unix(), task.GroupId).Scan(&task.Id)
	if err != nil {
		return models.Task{}, err
	}
	task.CreationDatetime = time.Now().Unix()
	task.CreatorId = userId
	return task, nil
}

//update task
func (db *DB) UpdateTask(UpdateTask models.Task, taskId int, userId int) (models.Task, error) {
	//
	task,_, err := db.GetTaskById(taskId)
	if err != nil {
		return models.Task{}, err
	}
	//
	if task.CreatorId != userId {
		return models.Task{}, fmt.Errorf("Id dont match ")
	}
	//
	_, err = db.Exec("UPDATE task_table SET assignee_id = $1, title = $2, description = $3, state = $4, deadline = $5," +
		" priority = $6 where id = $7",
		UpdateTask.AssigneeId, UpdateTask.Title, UpdateTask.Description, UpdateTask.State,
		UpdateTask.Deadline, UpdateTask.Priority, taskId)
	if err != nil {
		return models.Task{}, err
	}
	task.Title = UpdateTask.Title
	task.Description = UpdateTask.Description
	task.State = UpdateTask.State
	task.Deadline = UpdateTask.Deadline
	task.Priority = UpdateTask.Priority
	task.AssigneeId = UpdateTask.AssigneeId
	return task,nil
}
