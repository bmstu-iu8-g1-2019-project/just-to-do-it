package services

import "github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"

type DatastoreGroup interface {
	CreateGroup(models.Group) (error)
	GetGroup(int) (models.Group, error)
	UpdateGroup(int, models.Group) (error)
	DeleteGroup(int) (error)
}

//здесь вроде все примитивно, просто не знаю что еще можно накрутить по функционалу
func (db *DB) CreateGroup (group models.Group) (err error) {
	_, err = db.Exec("INSERT INTO group_table (title, description) values ($1, $2)",
		                    group.Title, group.Description)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetGroup(id int) (group models.Group, err error) {
	row := db.QueryRow("SELECT id, title, description FROM group_table WHERE id = $1", id)
	err = row.Scan(&group.Id, &group.Title, &group.Description)
	if err != nil {
		return group, err
	}
	return group, nil
}

func (db *DB) UpdateGroup(id int, group models.Group) (err error) {
	_, err = db.Exec("UPDATE group_table SET title = $1, description = $2 where id = $3", group.Title, group.Description, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteGroup (id int) (err error) {
	group, err := db.GetGroup(id)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM group_table WHERE id = $1", group.Id)
	if err != nil {
		return err
	}
	return nil
}
