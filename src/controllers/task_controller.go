package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"../services"
)

func GetTaskByTaskIdHandler(w http.ResponseWriter, r* http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	taskId, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	db := services.DB{}
	task := db.GetTaskByTaskId(taskId)

	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func GetTaskByAssigneeIdIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	assigneeId, _ := strconv.ParseInt(mux.Vars(r)["assignee_id"], 10, 64)

	db := services.DB{}
	tasks := db.GetTaskByAssigneeId(assigneeId)

	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func GetTaskByGroupIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	groupId, _ := strconv.ParseInt(mux.Vars(r)["assignee_id"], 10, 64)

	db := services.DB{}
	tasks := db.GetTaskByGroupId(groupId)

	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
