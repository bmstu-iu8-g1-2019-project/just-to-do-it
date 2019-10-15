package main

import (
	"log"
	"net/http"

	"dev-s/src/controllers"
	"dev-s/src/services"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

)

func main() {
	db, err := services.NewDB(services.OpenConfigFile("config.txt"))
	if err != nil {
		log.Panic(err)
	}


	env := &controllers.EnvironmentUser{ db}

	r := mux.NewRouter()
	r.HandleFunc("/login", env.ResponseLoginHandler).Methods("GET")
	r.HandleFunc("/register", env.ResponseRegisterHandler).Methods("POST")
	r.HandleFunc("/confirm", env.ConfirmEmailHandler).Methods("GET")
	r.HandleFunc("/user/{id}", env.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", env.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", env.DeleteUserHandler).Methods("DELETE")
	http.ListenAndServe(":3000", r)
}