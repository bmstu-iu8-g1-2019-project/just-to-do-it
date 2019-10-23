package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/gorilla/mux"
)

type EnvironmentTask struct {
	Db services.DatastoreTask
}

func(env *EnvironmentTask) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	assigneeIdStr := r.URL.Query().Get("assignee_id")
	title := r.URL.Query().Get("title")
	groupIdStr := r.URL.Query().Get("group_id")

	id := 0
	if idStr != "" {
		tmp, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id = tmp
	}

	assigneeId := 0
	if assigneeIdStr != "" {
		tmp, err := strconv.Atoi(assigneeIdStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		assigneeId = tmp
	}

	groupId := 0
	if groupIdStr != "" {
		tmp, err := strconv.Atoi(groupIdStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		groupId = tmp
	}

	tasks, err := env.Db.GetTask(id, assigneeId, title, groupId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(tasks)
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentTask)UpdateTask(w http.ResponseWriter, r *http.Request) {
	tmp := mux.Vars(r)["id"]
	if tmp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(tmp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task := models.Task{}
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = env.Db.UpdateTask(task, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentTask)CreateTask(w http.ResponseWriter, r *http.Request) {
	tmp := mux.Vars(r)["id"]
	if tmp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	assigneeId, err := strconv.Atoi(tmp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task := models.Task{}
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = env.Db.CreateTask(task, assigneeId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentTask) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	tmp := mux.Vars(r)["id"]
	if tmp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(tmp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.Db.DeleteTask(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
