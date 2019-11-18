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

// handle for user authorization
func (env *EnvironmentUser) ResponseLoginHandler(w http.ResponseWriter, r *http.Request) {
	receivedObject := models.User{}
	// write from the received data to the User structure
	err := json.NewDecoder(r.Body).Decode(&receivedObject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Message{Err: err, Msg: "Invalid parameters"})
		return
	}
	// function checks login and password in database
	user, err := env.Db.Login(receivedObject.Login, receivedObject.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_  = json.NewEncoder(w).Encode(models.Message{Err: err, Msg: "Not Found User"})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(models.Message{Err: err, Msg: "Wrong login or password"})
		return
	}
	user.Password = ""
	_ = json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

// user registration handle
func (env *EnvironmentUser) ResponseRegisterHandler (w http.ResponseWriter, r *http.Request) {
	obj := &models.User{}
	err := json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Message{Err: err, Msg: "Invalid parameters"})
		return
	}
	if obj.Password == ""  || obj.Login == "" || obj.Email == ""{
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Invalid parameters"})
		return
	}
	//function that detects the user and hashes his password
	user := models.User{}
	user, err = env.Db.Register(*obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Internal Server Error"})
		return
	}
	user.Password = ""
	_ = json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

// mail confirmation handle
func (env *EnvironmentUser) ConfirmEmailHandler (w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	err := env.Db.Confirm(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Internal Server Error"})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// user update handle
func (env *EnvironmentUser) UpdateUserHandler (w http.ResponseWriter, r *http.Request) {
	receivedObject := models.User{}
	// get the structure from the request
	err := json.NewDecoder(r.Body).Decode(&receivedObject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Invalid parameters"})
		return
	}
	// get from url id
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Invalid parameters"})
		return
	}
	// function updates by id and data from the request
	err = env.Db.UpdateUser(int(id), receivedObject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Internal Server Error"})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) GetUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Invalid parameters"})
		return
	}
	user, err := env.Db.GetUser(int(id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Not found user"})
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Internal server error"})
		return
	}
	_ = json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Invalid parameters"})
		return
	}
	err = env.Db.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Not found user"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.Message{Err : err, Msg: "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Message{Err : nil, Msg: "Deleted User"})
}
