package models

type Timetable struct {
	Id            int   `json:"id"`
	CreatorId     int   `json:"creator_id"`
	GroupId       int   `json:"group_id"`
	BeginInterval int64 `json:"begin_interval"`
	EndInterval   int64 `json:"end_interval"`
}
