package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
)

type EnvironmentTask struct {
	Db services.DatastoreTask
}

func(env *EnvironmentTask) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	strSlice := []string{}
	idStr := r.URL.Query().Get("id")
	assigneeIdStr := r.URL.Query().Get("assignee_id")
	groupIdStr := r.URL.Query().Get("group_id")
	strSlice = append(strSlice, idStr, assigneeIdStr, groupIdStr)
	title := r.URL.Query().Get("title")
	idSlice := []int{}

	for _, k := range strSlice {
		if k != "" {
			tmp, err := strconv.Atoi(k)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(err.Error())
				return
			}
			idSlice = append(idSlice, tmp)
		} else {
			idSlice = append(idSlice, 0)
		}
	}

	tasks, err := env.Db.GetTasks(idSlice, title)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(tasks)
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentTask) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	assigneeIdStr := r.URL.Query().Get("assignee_id")
	title := r.URL.Query().Get("title")
	groupIdStr := r.URL.Query().Get("group_id")

	taskId := 0
	if idStr != "" {
		tmp, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		taskId = tmp
	}

	assigneeId := 0
	if assigneeIdStr != "" {
		tmp, err := strconv.Atoi(assigneeIdStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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

	task, err := env.Db.GetTask(taskId, assigneeId, title, groupId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(task)
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentTask)UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id := 0
	if idStr != "" {
		tmp, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
		id = tmp
	}

	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = env.Db.UpdateTask(task, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentTask)CreateTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = env.Db.CreateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentTask) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id := 0
	if idStr != "" {
		tmp, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
		id = tmp
	}

	err := env.Db.DeleteTask(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
