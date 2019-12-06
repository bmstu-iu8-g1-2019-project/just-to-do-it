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
}

func (db *DB)CreateScope(scope models.Scope) (models.Scope, error) {
	// Получение интервала для которого insert_begin пересекает область
	// (Проверка begin_interval)
	result, err := db.Exec("SELECT id, creator_id, group_id, begin_interval, end_interval FROM " +
		"timetable WHERE begin_interval < $1 AND end_interval > $1", scope.BeginInterval)
	if err != nil {
		return models.Scope{}, err
	}
	min, _ := result.RowsAffected()
	if min > 0 {
		return models.Scope{}, fmt.Errorf("Invalid interval ")
	}
	// Проверка end_interval
	result, err = db.Exec("SELECT id, creator_id, group_id, begin_interval, end_interval FROM " +
		"timetable WHERE begin_interval < $1 AND end_interval > $1", scope.EndInterval)
	if err != nil {
		return models.Scope{}, err
	}
	min, _ = result.RowsAffected()
	if min > 0 {
		return models.Scope{}, fmt.Errorf("Invalid interval ")
	}
	// В случае если интервал не препятствует другим то добавляем запись в бд
	err = db.QueryRow("INSERT INTO timetable (creator_id, group_id, begin_interval, end_interval)" +
		"values ($1, $2, $3, $4) RETURNING id", scope.CreatorId, scope.GroupId,
		scope.BeginInterval, scope.EndInterval).Scan(&scope.Id)
	if err != nil {
		return models.Scope{}, err
	}
	return scope, nil
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
	query := "SELECT id, creator_id, group_id, begin_interval, end_interval FROM timetable WHERE "
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
	_, err := db.Exec("UPDATE timetable SET group_id = $1, begin_interval = $2," +
		"end_interval = $3 where id = $4", scopeId)
	if err != nil {
		return models.Scope{}, err
	}
	scope.Id = scopeId
	return scope, nil
}

func (db *DB)DeleteScope(scopeId int) (err error) {
	_, err = db.Exec("DELETE FROM timetable WHERE id = $1", scopeId)
	return err
}

func (db *DB)GetScope(scopeId int) (scope models.Scope, err error) {
	row := db.QueryRow("SELECT id, creator_id, group_id, begin_interval, end_interval FROM " +
		"timetable WHERE id = $1", scopeId)
	err = row.Scan(&scope.Id, &scope.CreatorId, &scope.GroupId,
		&scope.BeginInterval, &scope.EndInterval)
	if err != nil {
		return models.Scope{}, err
	}
	return scope, nil
}
