package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)


type DatastoreScope interface {
	GetScopes([]int) ([]models.Scope, error)
	UpdateScope(int, models.Scope) (models.Scope, error)
	DeleteScope(int) error
	CreateScope(models.Scope) (models.Scope, error)
	GetScope(int) (models.Scope, error)
	GetScopesWithInterval(int, int) ([]models.Scope, error)
	GetTasksFromScope(scopeId int) (tasks []models.Task, err error)

	AddTaskInScope(scopeId int, taskId int) (timetable models.Timetable, err error)
}

func (db *DB)CreateScope(scope models.Scope) (models.Scope, error) {
	// Получение интервала для которого insert_begin пересекает область
	// (Проверка begin_interval)
	result, err := db.Exec("SELECT id, creator_id, group_id, begin_interval, end_interval FROM " +
		"scope WHERE begin_interval < $1 AND end_interval > $1", scope.BeginInterval)
	if err != nil {
		return models.Scope{}, err
	}
	min, _ := result.RowsAffected()
	if min > 0 {
		return models.Scope{}, fmt.Errorf("Invalid interval ")
	}
	// Проверка end_interval
	result, err = db.Exec("SELECT id, creator_id, group_id, begin_interval, end_interval FROM " +
		"scope WHERE begin_interval < $1 AND end_interval > $1", scope.EndInterval)
	if err != nil {
		return models.Scope{}, err
	}
	min, _ = result.RowsAffected()
	if min > 0 {
		return models.Scope{}, fmt.Errorf("Invalid interval ")
	}
	// В случае если интервал не препятствует другим то добавляем запись в бд
	err = db.QueryRow("INSERT INTO scope (creator_id, group_id, begin_interval, end_interval)" +
		"values ($1, $2, $3, $4) RETURNING id", scope.CreatorId, scope.GroupId,
		scope.BeginInterval, scope.EndInterval).Scan(&scope.Id)
	if err != nil {
		return models.Scope{}, err
	}
	return scope, nil
}

func (db *DB)GetScopesWithInterval(begin int, end int) (scopes []models.Scope, err error){
	rows, err := db.Query("SELECT id, creator_id, group_id, begin_interval, end_interval FROM " +
		"scope WHERE begin_interval >= $1 AND end_interval <= $2", begin, end)
	if err != nil {
		return []models.Scope{}, err
	}
	for rows.Next() {
		scope := models.Scope{}
		err = rows.Scan(&scope.Id, &scope.CreatorId, &scope.GroupId, &scope.BeginInterval, &scope.EndInterval)
		if err != nil {
			return []models.Scope{}, err
		}
		scopes = append(scopes, scope)
	}
	return scopes, nil
}

func (db *DB)GetTasksFromScope(scopeId int) (tasks []models.Task, err error) {
	rows, err := db.Query("SELECT scope_id, task_id FROM timetable WHERE scope_id = $1", scopeId)
	if err != nil {
		return []models.Task{}, err
	}
	for rows.Next() {
		timetable := models.Timetable{}
		err = rows.Scan(&timetable.ScopeId, &timetable.TaskId)
		if err != nil {
			return []models.Task{}, err
		}
		task, _,err := db.GetTaskById(timetable.TaskId)
		if err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *DB)AddTaskInScope(scopeId int, taskId int) (timetable models.Timetable, err error) {
	// Получаем свободное время в scope'e
	freeTime , err := db.GetFreeTimeFromScope(scopeId)
	if err != nil {
		return models.Timetable{}, err
	}
	// Получаем инофрмацию о длительности задачи
	task, _, err := db.GetTaskById(taskId)
	if err != nil {
		return models.Timetable{}, err
	}
	// Сравниваем время выполнения задачи и свободное
	if freeTime < task.Duration {
		return models.Timetable{}, fmt.Errorf("Task time longer than scheduled ")
	}
	// Запись в вспомогательную бд
	_, err = db.Exec("INSERT INTO timetable (scope_id, task_id)" +
		" values ($1, $2)", scopeId, taskId)
	if err != nil {
		return models.Timetable{}, err
	}
	timetable.ScopeId = scopeId
	timetable.TaskId = taskId
	return timetable, nil
}

func (db *DB)GetFreeTimeFromScope(scopeId int) (freeTime int64, err error) {
	// Получаем все задачи из скоупа
	tasks, err := db.GetTasksFromScope(scopeId)
	if err != nil {
		return 0, err
	}
	// Получаем скоуп
	scope, err := db.GetScope(scopeId)
	if err != nil {
		return 0, err
	}
	// Время скоупа
	scopeDuration := scope.EndInterval - scope.BeginInterval
	// Сколько по времени занимают задачи в скоупе
	var busyTime int64
	for _, task := range tasks {
		busyTime += task.Duration
	}
	// Свободное время
	freeTime = scopeDuration - busyTime
	if freeTime < 0 {
		return 0, fmt.Errorf("Out of range ")
	}
	return freeTime, nil
}

func (db *DB)GetScopes(params []int) (scopes []models.Scope, err error) {
	// Формирование строки запроса
	// Создаем мапу для облегчения прохода по параметрам
	queryMap := make(map[string]interface{})
	if params[0] != 0 {
		queryMap["id"] = params[0]
	}
	if params[1] != 0 {
		queryMap["creator_id"] = params[1]
	}
	if params[2] != 0 {
		queryMap["group_id"] = params[2]
	}
	// Запрос без параметров
	query := "SELECT id, creator_id, group_id, begin_interval, end_interval FROM scope WHERE "
	// where == поле таблицы
	// value == значение поля
	var values []interface{}
	var where []string
	i := 1
	// Формирование параметров вида : "id = $i"
	for k, v := range queryMap {
		values = append(values, v)
		where = append(where, fmt.Sprintf("%s = $%s", k, strconv.Itoa(i)))
		i++
	}
	// Запрос
	rows, err := db.Query(query + strings.Join(where, " AND "), values...)
	if err != nil {
		return []models.Scope{}, err
	}
	// Сканирование результатов
	scopes = make([]models.Scope, 0)
	for rows.Next() {
		scope := models.Scope{}
		err = rows.Scan(&scope.Id, &scope.CreatorId, &scope.GroupId, &scope.BeginInterval, &scope.EndInterval)
		if err != nil {
			return []models.Scope{}, err
		}
		scopes = append(scopes, scope)
	}
	return scopes, nil
}

func (db *DB)UpdateScope(scopeId int, scope models.Scope) (models.Scope, error) {
	_, err := db.Exec("UPDATE scope SET group_id = $1, begin_interval = $2," +
		"end_interval = $3 where id = $4", scopeId)
	if err != nil {
		return models.Scope{}, err
	}
	scope.Id = scopeId
	return scope, nil
}

func (db *DB)DeleteScope(scopeId int) (err error) {
	_, err = db.Exec("DELETE FROM scope WHERE id = $1", scopeId)
	return err
}

func (db *DB)GetScope(scopeId int) (scope models.Scope, err error) {
	row := db.QueryRow("SELECT id, creator_id, group_id, begin_interval, end_interval FROM " +
		"scope WHERE id = $1", scopeId)
	err = row.Scan(&scope.Id, &scope.CreatorId, &scope.GroupId,
		&scope.BeginInterval, &scope.EndInterval)
	if err != nil {
		return models.Scope{}, err
	}
	return scope, nil
}
