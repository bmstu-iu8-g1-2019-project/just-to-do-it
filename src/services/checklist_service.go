package services

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

type DatastoreChecklist interface {
	CreateChecklist(models.Checklist) error
	GetChecklists(int) ([]models.Checklist, error)
	GetChecklist(int) (models.Checklist, error)
	UpdateChecklist(int, models.Checklist) error
	DeleteChecklist(int) error
}

//аналогично группам примитивный функционал
func(db *DB) CreateChecklist(checklist models.Checklist) error {
	_, err := db.Exec("INSERT INTO checklists_table (task_id, name, state)" +
		"values ($1, $2, $3)", checklist.TaskId, checklist.Name, checklist.State)
	if err != nil {
		return err
	}
	return nil
}

func(db *DB) GetChecklists(taskId int) (checklists []models.Checklist, err error) {
	rows, err := db.Query("SELECT * FROM checklists_table WHERE task_id = $1", taskId)
	if err != nil {
		return checklists, err
	}

	for rows.Next() {
		checklist := models.Checklist{} // & or no
		err = rows.Scan(&checklist.Id, &checklist.TaskId, &checklist.Name, &checklist.State)
		if err != nil {
			return checklists, err
		}
		checklists = append(checklists, checklist)
	}
	return checklists, nil
}

func(db *DB) GetChecklist(id int) (checklist models.Checklist, err error) {
	row := db.QueryRow("SELECT * FROM checklists_table WHERE id = $1", id)
	err = row.Scan(&checklist.Id, &checklist.TaskId, &checklist.Name, &checklist.State)
	if err != nil {
		return checklist, err
	}
	return checklist, nil
}

func(db *DB) UpdateChecklist(id int, checklist models.Checklist) error {
	_, err := db.Exec("UPDATE checklists_table SET task_id = $1, name = $2, state = $3",
		checklist.TaskId, checklist.Name, checklist.State)
	if err != nil {
		return err
	}
	return nil
}

func(db *DB) DeleteChecklist(id int) error {
	_, err := db.GetChecklist(id)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM checklists_table WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
