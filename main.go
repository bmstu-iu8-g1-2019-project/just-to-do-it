package main

import (
<<<<<<< HEAD
=======
	"github.com/gorilla/mux"
>>>>>>> d61f975f200c4f875988b87c1f23658a3979f24a
	"log"
	"os"
	"time"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/smtp"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

type Datastore interface {
	Login(string) (User, error)
	Register(r *http.Request) (error)
	Confirm(string) (error)
}

type Environment struct {
	db Datastore
}

type DB struct {
	*sql.DB
}
type User struct {
	Id int            `json:"id"`
	Email string      `json:"email"`
	Login string      `json:"login"`
	Fullname string   `json:"fullname"`
	Password string   `json:"password"`
	Acc_verified bool `json:"acc_verified"`
}

type Auth_confirmation struct {
	Login string `json:"login"`
	Hash string `json:"hash"`
	Deadline time.Time `json:"deadline"`
}

//login
func (db *DB) Login(login string) (obj User, err error) {
	row := db.QueryRow("SELECT * FROM usertab WHERE login = $1", login)
	err = row.Scan(&obj.Id, &obj.Email, &obj.Login, &obj.Fullname, &obj.Password, &obj.Acc_verified)
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

//register
func (db *DB) Register(r *http.Request) (err error) {
	obj := &User{}
	err = json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 8)

	_, err = db.Exec("INSERT INTO usertab (email, login, fullname, password, acc_verified) values ($1, $2, $3, $4, $5)",
		                    obj.Email, obj.Login, obj.Fullname, string(hashedPassword), obj.Acc_verified)
	if err != nil {
		return err
	}
	hashMail, err := sendMail(obj.Login, obj.Email)
	if err != nil {
		return err
	}
	err = db.recordMailConfirm(obj.Login, hashMail)
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
// Функция генерирует хэш, но пока без соли:(
func addressGenerator(login string) (str string) {
	h := sha256.New()
	h.Write([]byte(login))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Функция записывает в бд таблицу значения логина хэша и дедлайна
func (db *DB) recordMailConfirm (login string, hash string) (err error){
	deadlineTime := time.Now().Add(24 * time.Hour)
	_, err = db.Exec("INSERT INTO auth_confirmation (login, hash, deadline) values ($1, $2, $3)",
		                    login, hash, deadlineTime)
	if err != nil {
		return err
	}
	return nil
}

// Функция отправки сообщения на почту
func sendMail(login string, email string) (hash string, err error) {
	hash = addressGenerator(login)
	url := "\nlocalhost:3000/confirm?hash=" + hash
	from := os.Getenv("email")
	pass := os.Getenv("pass")

	msg := "\nFrom :" + from + "\n" +
		   "To: " + email + "\n" +
		   "Please confirm your email: " +
		   url

	err = smtp.SendMail("smtp.gmail.com:587",
		  smtp.PlainAuth(
		  	"",
		  	from,
		  	pass,
		  	"smtp.gmail.com"),
		  	from, []string{email}, []byte(msg))

	if err != nil {
		return "", err
	}
	return hash, nil
}


///confirm&hash=...
func (db *DB) Confirm(hash string) (err error) {
	var conf Auth_confirmation
	row := db.QueryRow("SELECT * FROM auth_confirmation WHERE hash = $1", hash)
	err = row.Scan(&conf.Login, &conf.Hash, &conf.Hash)
	_, err = db.Exec("UPDATE usertab SET acc_verified = $1 where login = $2", true, conf.Login)
	if err != nil {
		return err
	}
	return nil
}

func (env *Environment) confirmEmailHandler (w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	err := env.db.Confirm(hash)
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
	r.HandleFunc("/confirm", env.confirmEmailHandler).Methods("GET")
	http.ListenAndServe(":3000", r)
}




