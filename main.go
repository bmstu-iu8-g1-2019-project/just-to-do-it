package main

import (
	"authorization(krasivo)/authorization"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main () {
	db, err := authorization.NewDB("postgres://postgres:pass@localhost/postgres")
	if err != nil {
		log.Panic(err)
	}

	env := &authorization.Environment{db}

	r := mux.NewRouter()
	r.HandleFunc("/login", env.ResponseLoginHandler).Methods("GET")
	r.HandleFunc("/register", env.ResponseRegisterHandler).Methods("POST")
	r.HandleFunc("/confirm", env.ConfirmEmailHandler).Methods("GET")
	http.ListenAndServe(":3000", r)
}
