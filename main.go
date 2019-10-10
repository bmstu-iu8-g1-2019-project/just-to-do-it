package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/tree/dev-s/src/auth"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/tree/dev-s/src/db"
)

func main () {
	db, err := NewDB("postgres://postgres:pass@localhost/postgres")
	if err != nil {
		log.Panic(err)
	}

	env := &Environment{db}

	r := mux.NewRouter()
	r.HandleFunc("/login", env.responseLoginHandler).Methods("GET")
	r.HandleFunc("/register", env.responseRegisterHandler).Methods("POST")
	r.HandleFunc("/confirm", env.confirmEmailHandler).Methods("GET")
	r.HandleFunc("/user/{id}", env.updateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", env.getUserHandler).Methods("GET")
        r.HandleFunc("/user/{id}", env.deleteUserHandler).Methods("DELETE")
	http.ListenAndServe(":3000", r)
}
