package authorization

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

//login
func (db *DB) Login(login string) (obj User, err error) {
	rows := db.QueryRow("SELECT * FROM usertab WHERE login = $1", login)
	err = rows.Scan(&obj.Id, &obj.Email, &obj.Login, &obj.Fullname, &obj.Password, &obj.Acc_verified)
	if err != nil {
		return User{}, err
	}
	return obj, nil
}

func (env *Environment) ResponseLoginHandler(w http.ResponseWriter, r *http.Request) {
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

func (env *Environment) ResponseRegisterHandler(w http.ResponseWriter, r *http.Request) {
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

//confirm
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

func (env *Environment) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	err := env.db.Confirm(hash)
	if err != nil {
		return
	}
}
