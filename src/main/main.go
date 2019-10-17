package main

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	db, err := services.NewDB(services.OpenConfigFile("config.txt"))
	if err != nil {
		log.Panic(err)
	}

	env := &controllers.EnvironmentTask{ db}

	r := mux.NewRouter()
	r.Use(SetJSONHeader)

	r.HandleFunc("/user/task/{id}", env.GetTaskTIdHandler).Methods("GET")
	r.HandleFunc("/user/task/{assignee_id}", env.GetTasksAIdHandler).Methods("GET")
	r.HandleFunc("/user/task/{group_id}", env.GetTasksGIdHandler).Methods("GET")
	r.HandleFunc("/user/task", env.CreateTask).Methods("POST")
	r.HandleFunc("/user/task/{id}", env.UpdateTask).Methods("PUT")


}

func SetJSONHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}