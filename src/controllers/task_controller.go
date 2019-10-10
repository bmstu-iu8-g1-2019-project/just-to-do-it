package controllers

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func UpdateTaskAllProperties(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := services.DB{}
	err = db.UpdateTaskAllProperties(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := services.DB{}
	err = db.PostTask(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
