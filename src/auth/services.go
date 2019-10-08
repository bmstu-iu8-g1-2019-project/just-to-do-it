package auth

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

//login
func (db *DB) Login(login string) (obj User, err error) {
	row := db.QueryRow("SELECT * FROM usertab WHERE login = $1", login)
	err = row.Scan(&obj.Id, &obj.Email, &obj.Login, &obj.Fullname, &obj.Password, &obj.Acc_verified)
	if err != nil {
		return User{}, err
	}
	return obj, nil
}

//register
func (db *DB) Register(obj User) (err error) {
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


//confirm&hash=...
func (db *DB) Confirm(hash string) (err error) {
	var conf Auth_confirmation
	row := db.QueryRow("SELECT * FROM auth_confirmation WHERE hash = $1", hash)
	err = row.Scan(&conf.Login, &conf.Hash, &conf.Deadline)
	_, err = db.Exec("UPDATE usertab SET acc_verified = $1 where login = $2", true, conf.Login)
	if err != nil {
		return err
	}
	return nil
}
