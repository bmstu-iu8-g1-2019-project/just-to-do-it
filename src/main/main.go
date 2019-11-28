package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
)

func main() {
	db, err := services.NewDB()
	if err != nil {
		log.Panic(err)
	}

	envUser := &controllers.EnvironmentUser{Db: db}
	envTask := &controllers.EnvironmentTask{Db: db}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := mux.NewRouter()

	r.HandleFunc("/user/register", envUser.ResponseRegisterHandler).Methods("POST")
	r.HandleFunc("/user/login", envUser.ResponseLoginHandler).Methods("GET")
	r.HandleFunc("/user/{id}", envUser.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", envUser.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", envUser.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/confirm", envUser.ConfirmEmailHandler).Methods("GET")
	r.HandleFunc("/user/{id}/tasks", envTask.GetTasksHandler).Methods("GET")
	r.HandleFunc("/user/{id}/task/create", envTask.CreateTask).Methods("POST")

	err = http.ListenAndServe(":" + port, r)
	if err != nil {
		fmt.Println(err)
	}
}
