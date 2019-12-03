package services

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

type DatastoreGroup interface {
	CreateGroup(models.Group) (models.Group, error)
	GetGroup(int) (models.Group, error)
	UpdateGroup(int, models.Group) (models.Group, error)
	DeleteGroup(int) error
}

func (db *DB)CreateGroup(group models.Group) (models.Group, error) {
	err := db.QueryRow("INSERT INTO group_table (title, description) values ($1, $2)  RETURNING id",
		group.Title, group.Description).Scan(&group.Id)
	if err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func (db *DB)GetGroup(id int) (group models.Group, err error) {
	row := db.QueryRow("SELECT id, title, description FROM group_table WHERE id = $1", id)
	err = row.Scan(&group.Id, &group.Title, &group.Description)
	if err != nil {
		return group, err
	}
	return group, nil
}

func (db *DB)UpdateGroup(id int, group models.Group) (models.Group, error) {
	group.Id = id
	_, err := db.Exec("UPDATE group_table SET title = $1, description = $2 where id = $3", group.Title, group.Description, id)
	if err != nil {
		return group, err
	}
	return group, nil
}

func (db *DB)DeleteGroup(id int) (err error) {
	//Удаление всех задач из группы
	values := make([]int, 3)
	values[2] = id
	tasks, err := db.GetTasks(values, "", 0)
	for _, value := range tasks {
		_, err = db.Exec("DELETE FROM task_table WHERE id = $1", value.Id)
		if err != nil {
			return err
		}
	}
	//Удалние группы
	_, err = db.Exec("DELETE FROM group_table WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
