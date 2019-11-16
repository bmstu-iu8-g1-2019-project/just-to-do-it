package services

import (
	"fmt"
	"time"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"golang.org/x/crypto/bcrypt"
)

type DatastoreUser interface {
        Login(string, string) (models.User, error)
	Register(models.User) (models.User, error)
	Confirm(string) error
	UpdateUser(int, models.User) error
	GetUser(int) (models.User, error)
	DeleteUser(int) error
}

// The function checks the username and password that comes
// from the request with the username and password that lies in the database
func (db *DB)Login(login string, password string) (user models.User, err error) {
	row := db.QueryRow("SELECT * FROM user_table WHERE login = $1", login)
	err = row.Scan(&user.Id, &user.Email, &user.Login, &user.Fullname, &user.Password, &user.AccVerified)
	if err != nil {
		return user, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, err
	}
	return user,nil
}


// user registration (User structure comes)
func (db *DB) Register(obj models.User) (err error) {
	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 8)

	if obj.Fullname == "" {
		obj.Fullname = obj.Login
	}
	query := "INSERT INTO user_table (email, login, fullname, password, acc_verified)" +
		     "values ('%s', '%s', '%s', '%s', false)  RETURNING id"
	query = fmt.Sprintf(query, obj.Email, obj.Login, obj.Fullname, hashedPassword)

	err = db.QueryRow(query).Scan(&obj.Id)
	if err != nil {
		return obj, nil
	}
	//result, err := db.Exec("INSERT INTO user_table (email, login, fullname, password, acc_verified) values ($1, $2, $3, $4, $5)  RETURNING id",
	//	obj.Email, obj.Login, obj.Fullname, string(hashedPassword), false)
	//if err != nil {
	//	return obj, err
	//}
	//id, err := result.LastInsertId()
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println(id)
	//}
	//obj.Id = int(id)


	//write to the auxiliary table that stores the username its hash
	// and the expiration of the link to confirm mail
	err = db.recordMailConfirm(obj.Login)
	if err != nil {
		return obj, err
	}
	// send a message to the user's mail with such a login
	err = db.sendMail(obj.Login)
	if err != nil {
		return obj, err
	}
	return obj, err
}

// function that when clicking on the old link generates a new hash and resends the message to the mail
// else
// assigns "true" in the "acc_verified" field of the usera table and removes the entry from the auxiliary database
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

// update the user according to the parameters from the new request
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

// get the user structure by id
func (db *DB) GetUser (id int) (user models.User, err error) {
	row := db.QueryRow("SELECT * FROM user_table WHERE id = $1", id)
	err = row.Scan(&user.Id, &user.Email, &user.Login, &user.Fullname, &user.Password, &user.AccVerified)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// delete user by id
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
