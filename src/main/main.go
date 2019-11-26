package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
)

var (
	r = mux.NewRouter()
	Filename = "src/db/init.sql"

)

func init() {
	if err := godotenv.Load("vars.env"); err != nil {
		log.Print("No vars.env file found")
	}

	config := services.ReadConfig()
	fmt.Println(config)
	fmt.Println("Connecting to database server...")

	var db *services.DB
	chanDB := make(chan *services.DB, 1)
	timeout := time.After(time.Second * 16)
	go func() {
		for {
			db, err := services.NewDB(config)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Millisecond * 1500)
			} else {
				chanDB <- db
				return
			}
		}
	}()

MAIN:
	for {
		select {
		case database := <-chanDB:
			db = database
			log.Println("Connected to database")
			break MAIN
		case <-timeout:
			log.Println("Timout: connection was not established")
			panic("timeout")
		}
	}
	services.Setup(Filename, db)

	fmt.Println("URA!")
	envUser := &controllers.EnvironmentUser{db}
	envTask := &controllers.EnvironmentTask{db}


	r.Use(SetJSONHeader)

	// r.HandleFunc("/user/task/{id}", envTask.GetTaskTIdHandler).Methods("GET")
	// r.HandleFunc("/user/task/{assignee_id}", envTask.GetTasksAIdHandler).Methods("GET")
	// r.HandleFunc("/user/task/{group_id}", envTask.GetTasksGIdHandler).Methods("GET")
	r.HandleFunc("/user/task", envTask.CreateTask).Methods("POST")
	// r.HandleFunc("/user/task/{id}", envTask.UpdateTask).Methods("PUT")
	r.HandleFunc("/login", envUser.ResponseLoginHandler).Methods("GET")
	r.HandleFunc("/register", envUser.ResponseRegisterHandler).Methods("POST")
	r.HandleFunc("/confirm", envUser.ConfirmEmailHandler).Methods("GET")
	r.HandleFunc("/user/{id}", envUser.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", envUser.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", envUser.DeleteUserHandler).Methods("DELETE")
}

func main() {

	log.Fatal(http.ListenAndServe(":8080", r))
}

func SetJSONHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}
