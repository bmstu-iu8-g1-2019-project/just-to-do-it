package services

import (
	"fmt"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

func (db *DB)CreateLabel(label models.Label, taskId int) (models.Label, error) {
	query := "INSERT INTO label_table (task_id, title, color) values ('%d','%s','%s') RETURNING id"
	query = fmt.Sprintf(query, taskId, label.Title, label.Color)
	err := db.QueryRow(query).Scan(&label.Id)
	if err != nil {
		return models.Label{}, err
	}
	return label, nil
}

func (db *DB)GetLabelsByTaskId(taskId int) ([]models.Label, error) {
	rows, err := db.Query("SELECT * FROM label_table WHERE task_id = $1", taskId)
	if err != nil {
		return []models.Label{}, err
	}
	labels := make([]models.Label, 0)
	for rows.Next() {
		label := &models.Label{}
		err = rows.Scan(&label.Id, &label.TaskId, &label.Title, &label.Color)
		if err != nil {
			return []models.Label{}, err
		}
		labels = append(labels, *label)
	}
	return labels, nil
}

func (db *DB)GetLabel(labelId int) (label models.Label, err error) {
	row := db.QueryRow("SELECT * FROM label_table WHERE id = $1", labelId)
	err = row.Scan(&label.Id, &label.TaskId, &label.Title, &label.Color)
	if err != nil {
		return models.Label{}, err
	}
	return label, nil
}

func (db *DB)UpdateLabelColor(labelId int, color string) (label models.Label, err error) {
	label, err = db.GetLabel(labelId)
	if err != nil {
		return models.Label{}, err
	}
	_, err = db.Exec("UPDATE label_table SET color = $1 where id = $2", color, labelId)
	if err != nil {
		return models.Label{}, err
	}
	label.Color = color
	return label, nil
}

func (db *DB)UpdateLabelTitle(labelId int, title string) (label models.Label, err error) {
	label, err = db.GetLabel(labelId)
	if err != nil {
		return models.Label{}, err
	}
	_, err = db.Exec("UPDATE label_table SET title = $1 where id = $2", title, labelId)
	if err != nil {
		return models.Label{}, err
	}
	label.Title = title
	return label, nil
}

func (db *DB)DeleteLabel(labelId int) error {
	_, err := db.Exec("DELETE  FROM label_table WHERE id = $1", labelId)
	if err != nil {
		return err
	}
	return nil
}