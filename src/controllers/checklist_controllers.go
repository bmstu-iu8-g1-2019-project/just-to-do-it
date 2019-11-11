package controllers

import (
        "net/http"
	"strconv"
        "encoding/json"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/gorilla/mux"
)

type EnvironmentChecklist struct {
	Db services.DatastoreChecklist
}

func(env *EnvironmentChecklist) CreateChecklistHandler(w http.ResponseWriter, r *http.Request) {
	checklist := models.Checklist{}
	err := json.NewDecoder(r.Body).Decode(&checklist)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = env.Db.CreateChecklist(checklist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func(env *EnvironmentChecklist) GetChecklistHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	checklist, err := env.Db.GetChecklist(taskId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(checklist)
	w.WriteHeader(http.StatusOK)
}

func(env *EnvironmentChecklist) GetChecklistsHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	checklists, err := env.Db.GetChecklists(taskId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(checklists)
	w.WriteHeader(http.StatusOK)
}

func(env *EnvironmentChecklist) UpdateChecklistHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	checklist := models.Checklist{}
	err = json.NewDecoder(r.Body).Decode(&checklist)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = env.Db.UpdateChecklist(id, checklist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(checklist)
	w.WriteHeader(http.StatusOK)
}

func(env *EnvironmentChecklist) DeleteChecklistHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = env.Db.DeleteChecklist(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
