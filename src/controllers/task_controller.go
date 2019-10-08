package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"https://github.com/bmstu-iu8-g1-2019-project/just-to-do-it/tree/dev-d/src/services"
	"https://github.com/bmstu-iu8-g1-2019-project/just-to-do-it/tree/dev-d/src/models"
)

func TaskView(w http.ResponseWriter, r* http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	taskId, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	task := services.DB.GetTaskByTaskId(models.DB{}, taskId)

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
