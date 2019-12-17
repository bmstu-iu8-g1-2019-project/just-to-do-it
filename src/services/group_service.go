package services

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

type DatastoreGroup interface {
	CreateGroup(models.Group, int) (models.Group, error)
	GetGroup(int) (models.Group, error)
	UpdateGroup(int, models.Group) (models.Group, error)
	DeleteGroup(int) error
	GetGroups(userId int) ([]models.Group, error)
	//
	AddTaskInTrack(int, int) (models.TrackTaskPrevious, error)
	CreateTaskInTrack(int, int, models.Task) (models.TrackTaskPrevious, models.Task, error)
	DeleteTaskInTrack(int, int) error
	DeleteTrack(int) error
	UpdateTrack(int, models.Track) (models.Track, error)
	GetTrack(int) (models.Track, []models.Task, error)
	CreateTrack(int, models.Track) (models.Track, error)
}

func (db *DB)CreateGroup(group models.Group, userId int) (models.Group, error) {
	err := db.QueryRow("INSERT INTO group_table (creator_id, title, description) values ($1, $2, $3)  RETURNING id",
		userId, group.Title, group.Description).Scan(&group.Id)
	if err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func (db *DB)GetGroup(id int) (group models.Group, err error) {
	row := db.QueryRow("SELECT id, creator_id, title, description FROM group_table WHERE id = $1", id)
	err = row.Scan(&group.Id, &group.CreatorId, &group.Title, &group.Description)
	if err != nil {
		return group, err
	}
	return group, nil
}

func (db *DB)GetGroups(userId int) ([]models.Group, error) {
	groups := make([]models.Group, 0)
	rows, err := db.Query("SELECT id, creator_id, title, description FROM group_table WHERE creator_id = $1", userId)
	if err != nil {
		return []models.Group{}, err
	}
	for rows.Next() {
		group := models.Group{}
		err = rows.Scan(&group.Id, &group.CreatorId, &group.Title, &group.Description)
		if err != nil {
			return []models.Group{}, err
		}
		groups = append(groups, group)
	}
	return groups, nil
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
