package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"time"

	"dev-s/src/models"
	)

const (
	url = "\nlocalhost:3000/confirm?hash="
	from = "kolesnikov.school4@gmail.com"
	pass = "Proektk2019"
	msgConst = "\nFrom :%s\nTo: %s\nPlease confirm your email: %s"
)

func (db *DB) recordMailConfirm (login string) (err error){
	h := sha256.New()
	h.Write([]byte(login))
	hashed := h.Sum(nil)
	secret := hex.EncodeToString(hashed)

	deadlineTime := time.Now().Add(24 * time.Hour)
	_, err = db.Exec("INSERT INTO auth_confirmation (login, hash, deadline) values ($1, $2, $3)",
		login, string(secret), deadlineTime)
	if err != nil {
		return err
	}
	go db.resending(deadlineTime, login, string(secret))
	return nil
}

func (db *DB) resending(deadline time.Time, hash string, login string) (err error){
	timer := time.NewTimer(time.Until(deadline))
	<-timer.C
	row := db.QueryRow("SELECT * FROM user_table WHERE login = $1", login)
	var obj models.User
	err = row.Scan(&obj.Id, &obj.Email, &obj.Login, &obj.Fullname, &obj.Password, &obj.AccVerified)
	if err != nil {
		return err
	}
	if (obj.AccVerified == false) {
		err = db.sendMail(obj.Login, hash)
		if err != nil {
			return err
		}
		err = db.confirmFieldUpdate(login, hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) confirmFieldUpdate(login string, hash string) (err error) {
	_, err = db.Exec("UPDATE auth_confirmation SET hash = $1 where login = $2", hash, login)
	if err != nil {
		return err
	}
	return nil
}

// Mail sending function
func (db *DB) sendMail(login string, email string) (err error) {
	row := db.QueryRow("SELECT * FROM auth_confirmation WHERE login = $1", login)
	var obj models.AuthConfirmation
	err = row.Scan(&obj.Login, &obj.Hash, &obj.Deadline)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf(msgConst, from, email, url + obj.Hash)

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth(
			"",
			from,
			pass,
			"smtp.gmail.com"),
		from, []string{email}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}


