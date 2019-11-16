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
    db, err := services.NewDB()
    if err != nil {
	log.Panic(err)
    }

    envUser := &controllers.EnvironmentUser{ db}
    envTask := &controllers.EnvironmentTask{ db}
    envGroup := &controllers.EnvironmentGroup{db}

    r := mux.NewRouter()
    r.Use(SetJSONHeader)

    r.HandleFunc("/user/register", envUser.ResponseRegisterHandler).Methods("POST")
    r.HandleFunc("/user/login", envUser.ResponseLoginHandler).Methods("GET")
    r.HandleFunc("/user/{id}", envUser.UpdateUserHandler).Methods("PUT")
    r.HandleFunc("/user/{id}", envUser.GetUserHandler).Methods("GET")
    r.HandleFunc("/user/{id}", envUser.DeleteUserHandler).Methods("DELETE")
    r.HandleFunc("/confirm", envUser.ConfirmEmailHandler).Methods("GET")
    r.HandleFunc("/user/tasks/", envTask.GetTasksHandler).Methods("GET")
    r.HandleFunc("/user/task/", envTask.GetTaskHandler).Methods("GET")
    r.HandleFunc("/user/task/create", envTask.CreateTask).Methods("POST")
    r.HandleFunc("/user/task/", envTask.UpdateTask).Methods("PUT")
    r.HandleFunc("/user/task/", envTask.DeleteTaskHandler).Methods("DELETE")
    r.HandleFunc("/group", envGroup.CreateGroupHandler).Methods("POST")
    r.HandleFunc("/group/", envGroup.GetGroupHandler).Methods("GET")
    r.HandleFunc("/group/", envGroup.UpdateGroupHandler).Methods("PUT")
    r.HandleFunc("/group/", envGroup.DeleteGroupHandler).Methods("DELETE")
    http.ListenAndServe(":3000", r)
}

func SetJSONHeader(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	h.ServeHTTP(w, r)
    })
}
