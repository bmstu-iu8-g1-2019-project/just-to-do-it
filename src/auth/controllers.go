package auth

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (env *Environment) responseLoginHandler (w http.ResponseWriter, r *http.Request) {
	received_object := User{}
	err := json.NewDecoder(r.Body).Decode(&received_object)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	obj, err := env.db.Login(received_object.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(obj.Password), []byte(received_object.Password)); err != nil {
		http.Error(w, http.StatusText(401), 401)
	}
}

func (env *Environment) responseRegisterHandler (w http.ResponseWriter, r *http.Request) {
	obj := &User{}
	err := json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = env.db.Register(*obj)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		http.Error(w, http.StatusText(400), 400)
		return
	}
}

func (env *Environment) confirmEmailHandler (w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	err := env.db.Confirm(hash)
	if err != nil {
		return
	}
}
