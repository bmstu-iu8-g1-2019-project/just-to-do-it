package services

import (
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Datastore interface {
	Login(string) (models.User, error)
	Register(models.User) (error)
	Confirm(string) (error)
	UpdateUser(int, models.User) (error)
	GetUser(int) (models.User, error)
	DeleteUser(int) (error)
}

func (db *DB)Login(login string) (obj models.User, err error) {
	row := db.QueryRow("SELECT * FROM user_table WHERE login = $1", login)
	err = row.Scan(&obj.Id, &obj.Email, &obj.Login, &obj.Fullname, &obj.Password, &obj.AccVerified)
	if err != nil {
		return models.User{}, err
	}
	return obj, nil
}

func (db *DB) Register(obj models.User) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 8)

	_, err = db.Exec("INSERT INTO user_table (email, login, fullname, password, acc_verified) values ($1, $2, $3, $4, $5)",
		obj.Email, obj.Login, obj.Fullname, string(hashedPassword), obj.AccVerified)
	if err != nil {
		return err
	}
	err = db.recordMailConfirm(obj.Login)
	if err != nil {
		return err
	}
	err = db.sendMail(obj.Login)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Confirm(hash string) (err error) {
	var conf models.AuthConfirmation
	row := db.QueryRow("SELECT * FROM auth_confirmation WHERE hash = $1", hash)
	err = row.Scan(&conf.Login, &conf.Hash, &conf.Deadline)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if conf.Deadline.Before(time.Now()) {
		newHash := addressGenerator(conf.Login)
		err = db.confirmFieldUpdate(conf.Login, newHash)
		if err != nil {
			return err
		}
		err = db.sendMail(conf.Login)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	} else {
		_, err = db.Exec("UPDATE user_table SET acc_verified = true where login = $1", conf.Login)
		if err != nil {
			return err
		}
		_, err = db.Exec("DELETE  FROM auth_confirmation WHERE hash = $1", hash)
		if err != nil {
			return err
		}
		return nil
	}
}

func (db *DB) UpdateUser (id int, updateUser models.User) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 8)
	_, err = db.Exec("UPDATE user_table SET email = $1, login = $2, fullname = $3," +
		"password = $4, acc_verified = $5 where id = $6", updateUser.Email,
		updateUser.Login, updateUser.Fullname, hashedPassword, updateUser.AccVerified, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUser (id int) (user models.User, err error) {
	row := db.QueryRow("SELECT * FROM user_table WHERE id = $1", id)
	err = row.Scan(&user.Id, &user.Email, &user.Login, &user.Fullname, &user.Password, &user.AccVerified)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (db *DB) DeleteUser (id int) (err error) {
	user, err := db.GetUser(id)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE  FROM user_table WHERE id = $1", user.Id)
	if err != nil {
		return err
	}
	return nil
}
