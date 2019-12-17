package services

import (
	"database/sql"
	"fmt"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

func (db *DB)CreateTrack(groupId int, track models.Track) (models.Track, error) {
	track.GroupId = groupId
	err := db.QueryRow("INSERT INTO track_table (title, description, group_id) values ($1, $2, $3) RETURNING id",
		track.Title, track.Description, groupId).Scan(&track.Id)
	if err != nil {
		return models.Track{}, err
	}
	return track, nil
}

func (db *DB)GetTrack(id int) (track models.Track, tasks []models.Task, err error) {
	err = db.QueryRow("SELECT * FROM track_table WHERE id = $1", id).Scan(
		&track.Id, &track.Title, &track.Description, &track.GroupId)
	if err != nil {
		return models.Track{}, []models.Task{}, err
	}
	tasks, err = db.GetTasksInTrack(id)
	if err != nil {
		return models.Track{}, []models.Task{}, err
	}
	return track, tasks, nil
}

func (db *DB)GetTasksInTrack(id int) (tasks []models.Task, err error) {
	tasks = make([]models.Task, 0)
	rows, err := db.Query("SELECT * FROM track_task_previous WHERE track_id = $1", id)
	if err != nil {
		return []models.Task{}, err
	}
	for rows.Next() {
		var taskInTrack models.TrackTaskPrevious
		err = rows.Scan(&taskInTrack.TaskId, &taskInTrack.PreviousId, &taskInTrack.TrackId)
		if err != nil {
			return []models.Task{}, err
		}
		var task models.Task
		row := db.QueryRow("SELECT * FROM task_table WHERE id = $1", taskInTrack.TaskId)
		err = row.Scan(&task.Id, &task.CreatorId, &task.AssigneeId, &task.Title, &task.Description,
			&task.State, &task.Deadline, &task.Duration, &task.Priority, &task.CreationDatetime,
			&task.GroupId)
		if err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *DB)UpdateTrack(id int, track models.Track) (models.Track, error) {
	track.Id = id
	_, err := db.Exec("UPDATE track_table SET title = $1, description = $2, group_id = $3 WHERE id = $4",
		track.Title, track.Description, track.GroupId, id)
	if err != nil {
		return models.Track{}, nil
	}
	return track, nil
}

func (db *DB)DeleteTrack(id int) error {
	//удаление всех задач из трека
	rows, err := db.Query("SELECT * FROM track_task_previous WHERE track_id = $1", id)
	if err != nil {
		return err
	}

	for rows.Next() {
		taskFromTrack := models.TrackTaskPrevious{}
		err = rows.Scan(&taskFromTrack.TaskId, &taskFromTrack.PreviousId, &taskFromTrack.TrackId)
		if err != nil {
			return err
		}
		_, err = db.Exec("DELETE FROM task_table WHERE id = $1", taskFromTrack.TaskId)
		if err != nil {
			return err
		}
	}
	_, err = db.Exec("DELETE FROM track_task_previous WHERE track_id = $1", id)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM track_table WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB)AddTaskInTrack(taskId int, trackId int) (taskInTrack models.TrackTaskPrevious, err error) {
	// Проверка group_id
	task, _, err := db.GetTaskById(taskId)
	if err != nil {
		return models.TrackTaskPrevious{}, err
	}
	track, _, err := db.GetTrack(trackId)
	if err != nil {
		return models.TrackTaskPrevious{}, err
	}
	if task.GroupId != track.GroupId {
		return models.TrackTaskPrevious{}, fmt.Errorf("group_id don't match ")
	}
	// Проверка есть ли в таблице записи с таким trackId
	rows, err := db.Query("SELECT * FROM track_task_previous WHERE track_id = $1", trackId)
	if err != nil {
		return models.TrackTaskPrevious{}, err
	}
	for rows.Next() {
		err = rows.Scan(&taskInTrack.TaskId, &taskInTrack.PreviousId, &taskInTrack.TrackId)
	}
	// Если нет, то previous_id = 0
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO track_task_previous (task_id, previous_id, track_id)" +
			"values ($1, $2, $3)", taskId, 0, trackId)
		if err != nil {
			return models.TrackTaskPrevious{}, err
		}
		taskInTrack.TaskId = taskId
		taskInTrack.TrackId = trackId
		return taskInTrack, nil
	} else if err != nil {
		fmt.Println(err)
		return models.TrackTaskPrevious{}, err
	}
	// Иначе добавляем с previous_id = последняя добавленная задача
	_, err = db.Exec("INSERT INTO track_task_previous (task_id, previous_id, track_id)" +
		"values ($1, $2, $3)", taskId, taskInTrack.TaskId, trackId)
	if err != nil {
		return models.TrackTaskPrevious{}, err
	}
	taskInTrack.PreviousId = taskInTrack.TaskId
	taskInTrack.TrackId = trackId
	taskInTrack.TaskId = taskId
	return taskInTrack, nil
}

func (db *DB)CreateTaskInTrack(userId int, trackId int, task models.Task) (models.TrackTaskPrevious, models.Task, error) {
	track, _, err := db.GetTrack(trackId)
	if err != nil {
		return models.TrackTaskPrevious{}, models.Task{}, err
	}
	task.GroupId = track.GroupId
	task.CreatorId = userId
	task, err = db.CreateTask(task, userId)
	if err != nil {
		return models.TrackTaskPrevious{}, models.Task{}, err
	}
	taskInTrack, err := db.AddTaskInTrack(task.Id, trackId)
	if err != nil {
		return models.TrackTaskPrevious{}, models.Task{}, err
	}
	return taskInTrack, task, nil
}

func (db *DB)DeleteTaskInTrack(trackId int, taskId int) (err error) {
	deleteTask := models.TrackTaskPrevious{}
	err = db.QueryRow("SELECT * FROM track_task_previous WHERE task_id = $1",
		taskId).Scan(&deleteTask.TaskId, &deleteTask.PreviousId, &deleteTask.TrackId)
	if err != nil {
		return err
	}
	NextTaskAfterRemote := models.TrackTaskPrevious{}
	err = db.QueryRow("SELECT * FROM track_task_previous WHERE previous_id = $1", taskId).Scan(
		&NextTaskAfterRemote.TaskId, &NextTaskAfterRemote.PreviousId, &NextTaskAfterRemote.TrackId)
	if err == sql.ErrNoRows{
		_, err = db.Exec("DELETE FROM track_task_previous WHERE task_id = $1", taskId)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM track_task_previous WHERE task_id = $1", taskId)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM task_table WHERE id = $1", taskId)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE track_task_previous SET previous_id = $1 WHERE task_id = $2",
		deleteTask.PreviousId, NextTaskAfterRemote.TaskId)
	if err != nil {
		return err
	}

	return nil
}
