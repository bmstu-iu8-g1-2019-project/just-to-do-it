package main

import (
	"fmt"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var (
	FileName = "src/db/database.sql"
)

func init() {

	fmt.Println("Connecting to database server...")

	db, err := services.NewDB("postgres://docker:docker@localhost:5432/todoapp")
	if err != nil {
		fmt.Println(err)
	}
	envUser := &controllers.EnvironmentUser{ db}
	envTask := &controllers.EnvironmentTask{ db}

	r := mux.NewRouter()
	r.Use(SetJSONHeader)

	r.HandleFunc("/user/task/{id}", envTask.GetTaskTIdHandler).Methods("GET")
	r.HandleFunc("/user/task/{assignee_id}", envTask.GetTasksAIdHandler).Methods("GET")
	r.HandleFunc("/user/task/{group_id}", envTask.GetTasksGIdHandler).Methods("GET")
	r.HandleFunc("/user/task", envTask.CreateTask).Methods("POST")
	r.HandleFunc("/user/task/{id}", envTask.UpdateTask).Methods("PUT")
	r.HandleFunc("/login", envUser.ResponseLoginHandler).Methods("GET")
	r.HandleFunc("/register", envUser.ResponseRegisterHandler).Methods("POST")
	r.HandleFunc("/confirm", envUser.ConfirmEmailHandler).Methods("GET")
	r.HandleFunc("/user/{id}", envUser.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", envUser.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", envUser.DeleteUserHandler).Methods("DELETE")
}

func main() {

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetJSONHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}
