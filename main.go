package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
    "golang.org/x/crypto/bcrypt"
)

type Datastore interface {
	LoginHandler(string) (User, error)
	RegisterHandler(r *http.Request) (error)
}

type Environment struct {
	db Datastore
}

type DB struct {
	*sql.DB
}

type User struct {
	id int            `json:"id"`
	email string      `json:"email"`
	login string      `json:"login"`
	fullname string   `json:"fullname"`
	password string   `json:"password"`
	acc_verified bool `json:"acc_verified"`
}

//SignIn
func (db *DB) LoginHandler(login string) (obj User, err error) {
	rows := db.QueryRow("SELECT * FROM objectuser WHERE login = $1", login)

	err = rows.Scan(&obj.login, &obj.password)
	if err != nil {
		return User{}, err
	}

	return obj, nil
}

func (env *Environment) responseLoginHandler (w http.ResponseWriter, r *http.Request) {
	received_object := User{}
	err := json.NewDecoder(r.Body).Decode(received_object)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	obj, err := env.db.LoginHandler(received_object.login)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(obj.password), []byte(received_object.password)); err != nil {
		http.Error(w, http.StatusText(401), 401)
	}
}

//SignUp
func (db *DB) RegisterHandler(r *http.Request) (err error) {
	obj := User{}
	err = json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.password), 8)

	_, err = db.Exec("INSERT INTO objectuser (id, email, login, fullname, password, acc_verified) values ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING",
		                    obj.id, obj.email, obj.login, obj.fullname, string(hashedPassword), obj.acc_verified)
	if err != nil {
		return err
	}

	return nil
}

func (env *Environment) responseRegisterHandler (w http.ResponseWriter, r *http.Request) {
	err := env.db.RegisterHandler(r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
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

func main () {
	db, err := NewDB("postgres://postgres:pass@localhost/justtodoit")
	if err != nil {
		log.Panic(err)
	}

	env := &Environment{db}

	r := mux.NewRouter()
	r.HandleFunc("/login", env.responseLoginHandler).Methods("GET")
	r.HandleFunc("/register", env.responseRegisterHandler).Methods("POST")
	http.ListenAndServe(":3000", r)
}




