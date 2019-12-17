package models

import (
	"fmt"
)

type Scope struct {
	Id            int   `json:"id"`
	CreatorId     int   `json:"creator_id"`
	GroupId       int   `json:"group_id"`
	BeginInterval int64 `json:"begin_interval"`
	EndInterval   int64 `json:"end_interval"`
}

type Timetable struct {
	ScopeId int `json:"scope_id"`
	TaskId  int `json:"task_id"`
}

func ValidTimetable(scope Scope) (err error) {
	if  scope.CreatorId == 0 ||
		scope.GroupId == 0 ||
		scope.BeginInterval == -1 ||
		scope.EndInterval == -1 {
		return fmt.Errorf("Invalid body ")
	}
	if  scope.BeginInterval > scope.EndInterval ||
		scope.EndInterval < scope.BeginInterval {
		return fmt.Errorf("Invalid body ")
	}
	return nil
}
