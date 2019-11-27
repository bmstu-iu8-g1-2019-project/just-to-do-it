package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"

	_ "github.com/lib/pq"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	db, err := services.NewDB()
	if err != nil {
		log.Panic(err)
	}

	envUser := &controllers.EnvironmentUser{Db: db}
	envTask := &controllers.EnvironmentTask{Db: db}
	envGroup := &controllers.EnvironmentGroup{Db: db}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := mux.NewRouter()

	//user
	r.HandleFunc("/register", envUser.ResponseRegisterHandler).Methods("POST")
	r.HandleFunc("/login", envUser.ResponseLoginHandler).Methods("POST")
	r.HandleFunc("/user/{id}", envUser.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", envUser.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", envUser.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/confirm", envUser.ConfirmEmailHandler).Methods("GET")
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
	//group
	r.HandleFunc("/{id}/group/create", envGroup.CreateGroupHandler).Methods("POST")
	r.HandleFunc("/{id}/group/{group_id}/task/create", envTask.CreateTask).Methods("POST")
	r.HandleFunc("/{id}/group/{group_id}", envGroup.GetGroupHandler).Methods("GET")
	r.HandleFunc("/{id}/group/{group_id}", envGroup.UpdateGroupHandler).Methods("PUT")
	r.HandleFunc("/{id}/group/{group_id}", envGroup.DeleteGroupHandler).Methods("DELETE")
	//checklist
	r.HandleFunc("/{id}/task/{task_id}/checklist/create", envTask.CreateChecklistHandler).Methods("POST")
	r.HandleFunc("/{id}/checklist/{checklist_id}/item/create", envTask.CreateItemHandler).Methods("POST")
	r.HandleFunc("/{id}/checklist/{checklist_id}", envTask.GetChecklistHandler).Methods("GET")
	r.HandleFunc("/{id}/checklist/{checklist_id}", envTask.UpdateChecklistHandler).Methods("PUT")
	r.HandleFunc("/{id}/checklist/{checklist_id}", envTask.DeleteChecklistHandler).Methods("DELETE")
	r.HandleFunc("/{id}/checklist/{checklist_id}/items", envTask.GetChecklistItems).Methods("GET")
	r.HandleFunc("/{id}/checklist/{checklist_id}/item/{item_id}", envTask.UpdateItemHandler).Methods("PUT")
	r.HandleFunc("/{id}/item/{item_id}", envTask.DeleteItemHandler).Methods("DELETE")

	err = http.ListenAndServe(":" + port, r)
	if err != nil {
		fmt.Println(err)
	}
}
