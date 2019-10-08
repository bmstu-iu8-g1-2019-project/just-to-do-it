package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/smtp"
	"time"
)

//email verification
// Функция генерирует хэш
func addressGenerator(login string) (str string) {
	salt := make([]byte, 32)
	io.ReadFull(rand.Reader, salt)
	h := sha256.New()
	h.Write([]byte(login + string(salt)))
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
	go db.resending(deadlineTime, login)
	return nil
}

func (db *DB) resending(deadline time.Time, login string) {
	timer := time.NewTimer(time.Until(deadline))
	<-timer.C
	row := db.QueryRow("SELECT * FROM usertab WHERE login = $1", login)
	var obj User
	row.Scan(&obj.Id, &obj.Email, &obj.Login, &obj.Fullname, &obj.Password, &obj.Acc_verified)
	if (obj.Acc_verified == false) {
		hash, _ := sendMail(login, obj.Email)
		db.updateAccVerDB(login, hash)
	}
}

func (db *DB) updateAccVerDB(login string, hash string) {
	db.Exec("UPDATE auth_confirmation SET hash = $1 where login = $2", hash, login)
}

// Функция отправки сообщения на почту
func sendMail(login string, email string) (hash string, err error) {
	hash = addressGenerator(login)
	url := "\nlocalhost:3000/confirm?hash=" + hash
	from := "kolesnikov.school4@gmail.com"
	pass := "Proektk2019"

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
