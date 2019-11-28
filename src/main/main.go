package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

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
	envUser := &controllers.EnvironmentUser{Db: db}
	envGroup := &controllers.EnvironmentGroup{Db: db}

	//user
	r.HandleFunc("/register", envUser.ResponseRegisterHandler).Methods("POST")
	r.HandleFunc("/login", envUser.ResponseLoginHandler).Methods("POST")
	r.HandleFunc("/user/{id}", envUser.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", envUser.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", envUser.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/confirm", envUser.ConfirmEmailHandler).Methods("GET")
	//group
	r.HandleFunc("/{id}/group/create", envGroup.CreateGroupHandler).Methods("POST")
	r.HandleFunc("/{id}/group/{group_id}/task/create", envTask.CreateTask).Methods("POST")
	r.HandleFunc("/{id}/group/{group_id}", envGroup.GetGroupHandler).Methods("GET")
	r.HandleFunc("/{id}/group/{group_id}", envGroup.UpdateGroupHandler).Methods("PUT")
	r.HandleFunc("/{id}/group/{group_id}", envGroup.DeleteGroupHandler).Methods("DELETE")
}

func main() {

	log.Fatal(http.ListenAndServe(":8080", r))
}
