package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

func TaskView(w http.ResponseWriter, r* http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	taskId, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	task := (*services.DB).GetTaskByTaskId(&models.Environment{}, taskId)

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
