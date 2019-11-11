package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/gorilla/mux"
)

type EnvironmentUser struct {
	Db services.DatastoreUser
}

func (env *EnvironmentUser) ResponseLoginHandler(w http.ResponseWriter, r *http.Request) {
	received_object := models.User{}
	err := json.NewDecoder(r.Body).Decode(&received_object)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = env.Db.Login(received_object.Login, received_object.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) ResponseRegisterHandler (w http.ResponseWriter, r *http.Request) {
	obj := &models.User{}
	err := json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if obj.Password == ""  || obj.Login == "" || obj.Email == ""{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = env.Db.Register(*obj)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) ConfirmEmailHandler (w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	err := env.Db.Confirm(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) UpdateUserHandler (w http.ResponseWriter, r *http.Request) {
	received_object := models.User{}
	err := json.NewDecoder(r.Body).Decode(&received_object)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = env.Db.UpdateUser(int(id), received_object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) GetUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	received_object, err := env.Db.GetUser(int(id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(received_object)
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = env.Db.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
