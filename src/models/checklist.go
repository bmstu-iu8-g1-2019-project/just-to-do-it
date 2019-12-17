package models

type Checklist struct {
	Id     int    `json:"id"`
	TaskId int    `json:"task_id"`
	Name   string `json:"name"`
}

type ChecklistItem struct {
	Id          int    `json:"id"`
	ChecklistId int    `json:"checklist_id"`
	Name        string `json:"name"`
	State       string `json:"state"`
}
