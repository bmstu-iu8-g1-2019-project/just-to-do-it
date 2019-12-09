package models

type Track struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	GroupId     int    `json:"group_id"`
}

type TrackTaskPrevious struct {
	TaskId     int `json:"task_id"`
	PreviousId int `json:"previous_id"`
	TrackId    int `json:"track_id"`
}
