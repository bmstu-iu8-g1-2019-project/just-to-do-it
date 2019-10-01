package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Object struct {
	str string
}

type DB struct {
	*sql.DB
}

type Environment struct {
	db Datastore
}

type Datastore interface {
	GetData() ([]*Object, error)
	GetOneData(string) (Object, error)
}

func main () {
	db, err := NewDB("postgres://postgres:pass@localhost/serv")
	if err != nil {
		log.Panic(err)
	}

	env := &Environment{db}

	r := mux.NewRouter()
	r.HandleFunc("/serv", env.respAllData).Methods("GET")
	r.HandleFunc("/serv/{str}", env.respOneData).Methods("GET")
	http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}

// FOR /serv/{id}
func (db *DB) GetOneData(arg string) (Object, error) {
	rows := db.QueryRow("SELECT * FROM object WHERE str = $1", arg)

	var obj Object
	err := rows.Scan(&obj.str)
	if err != nil {
		return Object{}, err
	}
	return obj, nil
}

func (env *Environment) respOneData (w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["str"]
	obj, err := env.db.GetOneData(key)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	fmt.Fprint(w, obj.str)
}

// FOR /serv
func (db *DB) GetData() ([]*Object, error) {
	rows, err := db.Query("SELECT * FROM object")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	obj := make([]*Object, 0)
	for rows.Next() {
		ob := new(Object)
		err := rows.Scan(&ob.str)
		if err != nil {
			return nil, err
		}
		obj = append(obj, ob)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return obj, nil
}

func (env *Environment) respAllData(w http.ResponseWriter, r *http.Request) {
	obj, err := env.db.GetData()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, ob := range obj {
		fmt.Fprint(w, ob.str)
	}
}

// FOR DB
func NewDB(dbSourceName string) (*DB, error) {
	db, err := sql.Open("postgres",dbSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
