package main

import (
	"fmt"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var (
	FileName = "src/db/database.sql"

	Router = mux.NewRouter()
)

func init() {
	fmt.Println("Connecting to database server...")

	db, err := services.NewDB(os.Args[0])
	if err != nil {
		fmt.Println("Database opening error")
	}
	services.Setup(FileName, db)

	fmt.Println("Successfuly connection")

	envUser := &controllers.EnvironmentUser{ db}
	envTask := &controllers.EnvironmentTask{ db}


	Router.Use(SetJSONHeader)

	Router.HandleFunc("/user/task/{id}", envTask.GetTaskTIdHandler).Methods("GET")
	Router.HandleFunc("/user/task/{assignee_id}", envTask.GetTasksAIdHandler).Methods("GET")
	Router.HandleFunc("/user/task/{group_id}", envTask.GetTasksGIdHandler).Methods("GET")
	Router.HandleFunc("/user/task", envTask.CreateTask).Methods("POST")
	Router.HandleFunc("/user/task/{id}", envTask.UpdateTask).Methods("PUT")
	Router.HandleFunc("/login", envUser.ResponseLoginHandler).Methods("GET")
	Router.HandleFunc("/register", envUser.ResponseRegisterHandler).Methods("POST")
	Router.HandleFunc("/confirm", envUser.ConfirmEmailHandler).Methods("GET")
	Router.HandleFunc("/user/{id}", envUser.UpdateUserHandler).Methods("PUT")
	Router.HandleFunc("/user/{id}", envUser.GetUserHandler).Methods("GET")
	Router.HandleFunc("/user/{id}", envUser.DeleteUserHandler).Methods("DELETE")
}

func main() {
	log.Fatal(http.ListenAndServe(":5555", Router))
}

func SetJSONHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}
