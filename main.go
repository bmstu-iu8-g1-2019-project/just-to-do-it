package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/smtp"
)

type Datastore interface {
	Login(string) (User, error)
	Register(r *http.Request) (error)
	Confirm(r *http.Request) (error)
}

type Environment struct {
	db Datastore
}

type DB struct {
	*sql.DB
}

type User struct {
	id int            `json:"id"`
	Email string      `json:"email"`
	Login string      `json:"login"`
	Fullname string   `json:"fullname"`
	Password string   `json:"password"`
	Acc_verified bool `json:"acc_verified"`
}

//SignIn
func (db *DB) Login(login string) (obj User, err error) {
	rows := db.QueryRow("SELECT * FROM objectuser WHERE login = $1", login)
	err = rows.Scan(&obj.id, &obj.Email, &obj.Login, &obj.Fullname, &obj.Password, &obj.Acc_verified)
	if err != nil {
		return User{}, err
	}
	return obj, nil
}

func (env *Environment) responseLoginHandler (w http.ResponseWriter, r *http.Request) {
	received_object := User{}
	err := json.NewDecoder(r.Body).Decode(&received_object)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	obj, err := env.db.Login(received_object.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(obj.Password), []byte(received_object.Password)); err != nil {
		http.Error(w, http.StatusText(401), 401)
	}
}

//SignUp
func (db *DB) Register(r *http.Request) (err error) {
	obj := &User{}
	err = json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 8)

	_, err = db.Exec("INSERT INTO objectuser (email, login, fullname, password, acc_verified) values ($1, $2, $3, $4, $5)",
		                    obj.Email, obj.Login, obj.Fullname, string(hashedPassword), obj.Acc_verified)
	if err != nil {
		return err
	}
	url := "localhost:3000/confirm/" + obj.Login
	err = sendMail(url, obj.Email)
	if err != nil {
		return err
	}
	return nil
}

func (env *Environment) responseRegisterHandler (w http.ResponseWriter, r *http.Request) {
	err := env.db.Register(r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		http.Error(w, http.StatusText(400), 400)
		return
	}
}

//email verification
func addressGenerator() (str string) {

	return str
}

func sendMail(url string, email string) (err error) {
	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", "kolesnikov.school4@gmail.com", "Proektk2019", "smtp.gmail.com"),
		"kolesnikov.school4@gmail.com", []string{email}, []byte(url))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Confirm(r *http.Request) (err error) {
	params := mux.Vars(r)
	rows := db.QueryRow("SELECT * FROM objectuser WHERE login = $1", params["str"])
	obj := &User{}
	err = rows.Scan(&obj.id, &obj.Email, &obj.Login, &obj.Fullname, &obj.Password, &obj.Acc_verified)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE objectuser SET acc_verified = $1 where login = $2", true, obj.Login)
	if err != nil {
		return err
	}
	return nil
}

func (env *Environment) confirmEmailHandler (w http.ResponseWriter, r *http.Request) {
	err := env.db.Confirm(r)
	if err != nil {
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
	db, err := NewDB("postgres://postgres:pass@localhost/postgres")
	if err != nil {
		log.Panic(err)
	}

	env := &Environment{db}

	r := mux.NewRouter()
	r.HandleFunc("/login", env.responseLoginHandler).Methods("GET")
	r.HandleFunc("/register", env.responseRegisterHandler).Methods("POST")
	r.HandleFunc("/confirm/{str}", env.confirmEmailHandler).Methods("PUT")
	http.ListenAndServe(":3000", r)
}








