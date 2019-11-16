package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Tags struct {
	Id int       `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type Task struct {
	Id               int       `json:"id"`
	AssigneeId       int       `json:"assignee_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	State            string    `json:"state"`
	Deadline         time.Time `json:"deadline"`
	Priority         int       `json:"priority"`
	CreationDatetime time.Time `json:"creation_datetime"`
	GroupId          int       `json:"group_id"`
	Tag              []Tags    `json:"tags"`
}

func (t *Tags) Value() (driver.Value, error) {
	return fmt.Sprintf("(%d,'%s','%s')", t.Id, t.Title, t.Color), nil
}

