package services

import "github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"

type DatastoreTimeTable interface {
	GetTimetables([]int) (models.Timetable, error)
	UpdateTimetable(int, models.Timetable) (models.Timetable, error)
	DeleteTimetable(int) error
}
