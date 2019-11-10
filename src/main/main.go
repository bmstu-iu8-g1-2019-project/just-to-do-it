package main

import (
	"flag"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"log"
	"net/http"
	"fmt"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	dbHost := flag.String("db.host", "db", "db host")
	dbPort := flag.Int("db.port", 5432, "db port")
	dbUser := flag.String("db.user", "user", "db user")
	dbPassword := flag.String("db.password", "password", "db password")
	dbDatabase := flag.String("db.database", "db", "database name")
	flag.Parse()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", *dbHost, *dbPort, *dbUser, *dbPassword, *dbDatabase)
	db, err := services.NewDB(dsn)

	if err != nil {
        log.Panic(err)
    }

	// Not added migrations files yet
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			panic(err)
		}
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
    http.ListenAndServe(*httpAddr, r)
}

func SetJSONHeader(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        h.ServeHTTP(w, r)
    })
}
