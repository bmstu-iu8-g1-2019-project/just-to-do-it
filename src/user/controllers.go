package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (env *Environment) getUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	received_object, err := env.db.getUser(int(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(401), 401)
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(received_object)
}

func (env *Environment) updateUserHandler (w http.ResponseWriter, r *http.Request) {
	received_object := User{}
	err := json.NewDecoder(r.Body).Decode(&received_object)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	err = env.db.updateUser(int(id), received_object)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}