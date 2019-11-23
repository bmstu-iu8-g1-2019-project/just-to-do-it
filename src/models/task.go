package models

import (
	"time"
)

type Task struct {
	Id               int       `json:"id"`
	CreatorId        int       `json:"creator_id"`
	AssigneeId       int       `json:"assignee_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	State            string    `json:"state"`
	Deadline         time.Time `json:"deadline"`
	Priority         int       `json:"priority"`
	CreationDatetime time.Time `json:"creation_datetime"`
	GroupId          int       `json:"group_id"`
}

type Tags struct {
	Id     int    `json:"id"`
	TaskId int    `json:"task_id"`
	Title  string `json:"title"`
	Color  string `json:"color"`
}