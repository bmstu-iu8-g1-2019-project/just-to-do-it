package controllers

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)


func GetTaskTIdHandler(w http.ResponseWriter, r* http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	tmp := mux.Vars(r)["assignee_id"]
	if tmp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskId, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db := services.DB{}
	task := db.GetTaskTId(taskId)

	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func GetTasksAIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	tmp := mux.Vars(r)["assignee_id"]
	if tmp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	assigneeId, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	db := services.DB{}
	tasks := db.GetTasksAId(assigneeId)

	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func GetTasksGIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	tmp := mux.Vars(r)["assignee_id"]
	if tmp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupId, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db := services.DB{}
	tasks := db.GetTasksGId(groupId)

	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db := services.DB{}
	err = db.UpdateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db := services.DB{}
	err = db.CreateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
