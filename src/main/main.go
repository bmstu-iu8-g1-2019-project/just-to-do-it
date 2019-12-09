package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
)

var (
	r = mux.NewRouter()
	Filename = "src/db/init.sql"

)

func init() {
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
	envTask := &controllers.EnvironmentTask{Db: db}

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
	//task
	r.HandleFunc("/{id}/tasks", envTask.GetTasksHandler).Methods("GET")
	r.HandleFunc("/{id}/task/create", envTask.CreateTask).Methods("POST")
	r.HandleFunc("/{id}/task/{task_id}", envTask.GetTaskHandler).Methods("GET")
	r.HandleFunc("/{id}/task/{task_id}", envTask.UpdateTaskHandler).Methods("PUT")
	//label
	r.HandleFunc("/{id}/task/{task_id}/label/create", envTask.CreateLabelHandler).Methods("POST")
	r.HandleFunc("/{id}/label/{label_id}", envTask.GetLabelHandler).Methods("GET")
	r.HandleFunc("/{id}/task/{task_id}/labels", envTask.GetLabelsByTaskIdHandler).Methods("GET")
	r.HandleFunc("/{id}/label/{label_id}/color", envTask.UpdateLabelColorHandler).Methods("PUT")
	r.HandleFunc("/{id}/label/{label_id}/title", envTask.UpdateLabelTitleHandler).Methods("PUT")
	r.HandleFunc("/{id}/label/{label_id}", envTask.DeleteLabelHandler).Methods("DELETE")
	//checklist
	r.HandleFunc("/{id}/task/{task_id}/checklist/create", envTask.CreateChecklistHandler).Methods("POST")
	r.HandleFunc("/{id}/checklist/{checklist_id}/item/create", envTask.CreateItemHandler).Methods("POST")
	r.HandleFunc("/{id}/checklist/{checklist_id}", envTask.GetChecklistHandler).Methods("GET")
	r.HandleFunc("/{id}/checklist/{checklist_id}", envTask.UpdateChecklistHandler).Methods("PUT")
	r.HandleFunc("/{id}/checklist/{checklist_id}", envTask.DeleteChecklistHandler).Methods("DELETE")
	r.HandleFunc("/{id}/checklist/{checklist_id}/items", envTask.GetChecklistItems).Methods("GET")
	r.HandleFunc("/{id}/checklist/{checklist_id}/item/{item_id}", envTask.UpdateItemHandler).Methods("PUT")
	r.HandleFunc("/{id}/item/{item_id}", envTask.DeleteItemHandler).Methods("DELETE")
	//track
	r.HandleFunc("/{id}/track/{track_id}/task/create", envGroup.CreateTaskInTrackHandler).Methods("POST")
	r.HandleFunc("/{id}/track/{track_id}/task/{task_id}", envGroup.AddTaskInTrackHandler).Methods("POST")
	r.HandleFunc("/{id}/group/{group_id}/track/create", envGroup.CreateTrackHandler).Methods("POST")
	r.HandleFunc("/{id}/track/{track_id}", envGroup.GetTrackHandler).Methods("GET")
	r.HandleFunc("/{id}/track/{track_id}", envGroup.UpdateTrackHandler).Methods("PUT")
	r.HandleFunc("/{id}/track/{track_id}", envGroup.DeleteTrackHandler).Methods("DELETE")
	r.HandleFunc("/{id}/track/{track_id}/task/{task_id}", envGroup.DeleteTaskInTrack).Methods("DELETE")
}

func main() {

	log.Fatal(http.ListenAndServe(":8080", r))
}
