package models

type Checklist struct {
	Id int       `json:"id"`
	TaskId int   `json:"task_id"`
	Name string  `json:"name"`
	State string `json:"state"`
}
