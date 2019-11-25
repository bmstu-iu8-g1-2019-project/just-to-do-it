package services
//
//import (
//	"fmt"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
//)
//
//type DatastoreTask interface {
//	GetTasks([]int, string, int) ([]models.Task, error)
//	GetTaskById(int) (models.Task, error)
//	UpdateTask(models.Task, int) (models.Task, error)
//	CreateTask(models.Task, int) (models.Task, error)
//	//DeleteTask(int) error
//}
//
//// input we get an array of values from url
//// returns an array of objects of type Task
//func(db* DB) GetTasks(idSlice []int, title string, userId int) (tasks []models.Task ,err error) {
//	queryMap := make(map[string]interface{})
//	if idSlice[0] != 0 {
//		queryMap["id"] = idSlice[0]
//	}
//	if idSlice[1] != 0 {
//		queryMap["assignee_id"] = idSlice[1]
//	}
//	if idSlice[2] != 0 {
//		queryMap["group_id"] = idSlice[2]
//	}
//	if title != "" {
//		queryMap["title"] = title
//	}
//	if userId != 0 {
//		queryMap["creator_id"]= userId
//	}
//	query := "SELECT id, creator_id, assignee_id, title, description, state, deadline, priority, creation_datetime, group_id FROM task_table WHERE "
//
//	var values []interface{}
//	var where []string
//	i := 1
//	for k, v := range queryMap {
//		values = append(values, v)
//		where = append(where, fmt.Sprintf("%s = $%s", k, strconv.Itoa(i)))
//		i++
//	}
//	//fmt.Println(query + strings.Join(where, " AND "))
//
//	rows, err := db.Query(query + strings.Join(where, " AND "), values...)
//	if err != nil {
//		return []models.Task{}, err
//	}
//
//	tasks = make([]models.Task, 0)
//
//	for rows.Next() {
//		task := &models.Task{}
//		err = rows.Scan(&task.Id, &task.CreatorId, &task.AssigneeId, &task.Title, &task.Description,
//			&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime,
//			&task.GroupId)
//		if err != nil {
//			return tasks, err
//		}
//
//		tasks = append(tasks, *task)
//	}
//	return tasks, nil
//}
//
////get task
//func (db *DB) GetTaskById (id int) (task models.Task, err error) {
//	row := db.QueryRow("SELECT * FROM task_table WHERE id = $1", id)
//	err = row.Scan(&task.Id, &task.CreatorId, &task.AssigneeId, &task.Title, &task.Description,
//		&task.State, &task.Deadline, &task.Priority, &task.CreationDatetime, &task.GroupId)
//	if err != nil {
//		return models.Task{}, err
//	}
//	return task, nil
//}
//
////create task
//func (db *DB) CreateTask(task models.Task, userId int) (models.Task, error) {
//	query := "INSERT INTO task_table (creator_id, assignee_id, title, description, state, deadline, priority, creation_datetime, group_id) values ('%d', '%d', '%s', '%s', '%s', '%d', '%d', '%d', '%d')  RETURNING id"
//	query = fmt.Sprintf(query, userId, task.AssigneeId, task.Title, task.Description, task.State,
//		                task.Deadline, task.Priority, time.Now().Unix(), task.GroupId)
//	err := db.QueryRow(query).Scan(&task.Id)
//	if err != nil {
//		return models.Task{}, err
//	}
//	task.CreationDatetime = time.Now().Unix()
//	task.CreatorId = userId
//	return task,nil
//}
//
////update task
//func (db *DB) UpdateTask(UpdateTask models.Task, Id int) (models.Task, error) {
//	task, err := db.GetTaskById(Id)
//	_, err = db.Exec("UPDATE task_table SET assignee_id = $1, title = $2, description = $3, state = $4, deadline = $5," +
//		" priority = $6 where id = $7",
//		UpdateTask.AssigneeId, UpdateTask.Title, UpdateTask.Description, UpdateTask.State,
//		UpdateTask.Deadline, UpdateTask.Priority, Id)
//	if err != nil {
//		return models.Task{}, err
//	}
//	task.AssigneeId = UpdateTask.AssigneeId
//	task.Title = UpdateTask.Title
//	task.Description = UpdateTask.Description
//	task.State = UpdateTask.State
//	task.Deadline = UpdateTask.Deadline
//	task.Priority = UpdateTask.Priority
//	return task,nil
//}
