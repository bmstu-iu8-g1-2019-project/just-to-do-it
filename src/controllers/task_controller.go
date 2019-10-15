package controllers

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type EnvironmentUser struct {
	Db services.Datastore
}

func (env *EnvironmentUser)GetTaskTIdHandler(w http.ResponseWriter, r* http.Request) {
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

	task := env.Db.GetTaskTId(taskId)

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

func (env *EnvironmentUser)GetTasksAIdHandler(w http.ResponseWriter, r *http.Request) {
	tmp := mux.Vars(r)["assignee_id"]
	if tmp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	assigneeId, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	tasks := env.Db.GetTasksAId(assigneeId)

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

func (env *EnvironmentUser)GetTasksGIdHandler(w http.ResponseWriter, r *http.Request) {
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

	tasks := env.Db.GetTasksGId(groupId)

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

func (env *EnvironmentUser)UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = env.Db.UpdateTask(task)
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

func (env *EnvironmentUser)CreateTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = env.Db.CreateTask(task)
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
