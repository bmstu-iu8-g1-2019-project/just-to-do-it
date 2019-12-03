package models

type Task struct {
	Id               int       `json:"id"`
	CreatorId        int       `json:"creator_id"`
	AssigneeId       int       `json:"assignee_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	State            string    `json:"state"`
	Deadline         int64     `json:"deadline"`
	Duration         int64     `json:"duration"`
	Priority         int       `json:"priority"`
	CreationDatetime int64     `json:"creation_datetime"`
	GroupId          int       `json:"group_id"`
}

type Label struct {
	Id     int    `json:"id"`
	TaskId int    `json:"task_id"`
	Title  string `json:"title"`
	Color  string `json:"color"`
}
