package services

import(
	"fmt"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

func (db *DB)CreateChecklist(checklist models.Checklist, taskId int) (models.Checklist, error) {
	err := db.QueryRow("INSERT INTO checklist_table (task_id, name) values ($1, $2)  RETURNING id",
		taskId, checklist.Name).Scan(&checklist.Id)
	if err != nil {
		return models.Checklist{}, err
	}
	checklist.TaskId = taskId
	return checklist, nil
}

func (db *DB)CreateChecklistItem(item models.ChecklistItem, checklistId int) (models.ChecklistItem, error) {
	err := db.QueryRow("INSERT INTO checklistItem_table (checklist_id, name, state) values ($1, $2, $3)  RETURNING id",
		checklistId, item.Name, item.State).Scan(&item.Id)
	if err != nil {
		return models.ChecklistItem{}, err
	}
	item.ChecklistId = checklistId
	return item, nil
}

func (db *DB)GetChecklist(checklistId int) (models.Checklist, []models.ChecklistItem, error) {
	checklist := models.Checklist{}
	row := db.QueryRow("SELECT * FROM checklist_table WHERE id = $1", checklistId)
	err := row.Scan(&checklist.Id, &checklist.TaskId, &checklist.Name)
	if err != nil {
		return models.Checklist{}, []models.ChecklistItem{}, err
	}
	items, err := db.GetChecklistItems(checklist.Id)
	if err != nil {
		return models.Checklist{}, []models.ChecklistItem{}, err
	}
	return checklist, items, nil
}

func (db *DB)UpdateChecklist(checklistId int, checklist models.Checklist) (models.Checklist, error) {
	listInDb,_, err := db.GetChecklist(checklistId)
	if err != nil {
		return models.Checklist{}, err
	}
	_, err = db.Exec("UPDATE checklist_table SET name = $1 where id = $2", checklist.Name, checklistId)
	if err != nil {
		return models.Checklist{}, err
	}
	checklist.Id = checklistId
	checklist.TaskId = listInDb.TaskId
	return checklist, nil
}

func (db *DB)DeleteChecklist(checklistId int) error {
	items, err := db.GetChecklistItems(checklistId)
	if err != nil {
		return err
	}
	for _, value := range items {
		_, err = db.Exec("DELETE FROM checklistItem_table where id = $1", value.Id)
		if err != nil {
			return err
		}
	}
	_, err = db.Exec("DELETE FROM checklist_table where id = $1", checklistId)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB)GetChecklistItems(checklistId int) ([]models.ChecklistItem, error) {
	rows, err := db.Query("SELECT * FROM checklistItem_table " +
		"WHERE checklist_id = $1", checklistId)
	if err != nil {
		return []models.ChecklistItem{}, err
	}

	items := make([]models.ChecklistItem, 0)
	for rows.Next() {
		item := models.ChecklistItem{}
		err = rows.Scan(&item.Id, &item.ChecklistId,
			&item.Name, &item.State)
		if err != nil {
			return []models.ChecklistItem{}, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (db *DB)GetChecklistItem(itemId int) (models.ChecklistItem, error) {
	item := models.ChecklistItem{}
	row := db.QueryRow("SELECT * FROM checklistItem_table where id = $1", itemId)
	err := row.Scan(&item.Id, &item.ChecklistId, &item.Name, &item.State)
	if err != nil {
		return models.ChecklistItem{}, err
	}
	return item, nil
}

func (db *DB)UpdateChecklistItem(itemId int, checklistId int, item models.ChecklistItem) (models.ChecklistItem, error) {
	_, err := db.Exec("UPDATE checklistItem_table SET name = $1, state = $2 where id = $3",
		item.Name, item.State, itemId)
	if err != nil {
		return models.ChecklistItem{}, err
	}
	item.ChecklistId = checklistId
	item.Id = itemId
	return item, nil
}

func (db *DB)DeleteChecklistItem(itemId int) error {
	_, err := db.Exec("DELETE FROM checklistItem_table where id = $1", itemId)
	if err != nil {
		return err
	}
	return nil
}
