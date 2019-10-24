package tests

import (
	"../src/controllers"
	"../src/models"
	"../src/services"
	"golang.org/x/crypto/bcrypt"
	"gotest.tools/assert"
	"testing"
)

type mockDB struct {
	DB services.DatastoreUser
}

var (
	expectedUser1 = models.User{
		Id:	1,
		Email: "dk@mail.ru",
		Login: "dk",
		Fullname: "Dan Kokin",
		Password: "pass",
		AccVerified: true,
	}

	expectedUser2 = models.User{
		Id:	2,
		Email: "dk@mail.ru",
		Login: "dk",
		Fullname: "Dan Kokin",
		Password: "pass",
		AccVerified: true,
	}

	expectedUser3 = models.User{
		Id:	42,
		Email: "dk@mail.ru",
		Login: "dk",
		Fullname: "Dan Kokin",
		Password: "pass",
		AccVerified: true,
	}
)



func TestGetUser (t *testing.T) {
	users := make([]models.User, 0)
	users = append(users, expectedUser1)
	users = append(users, expectedUser2)
	users = append(users, expectedUser3)

	dbTest, _ := services.NewMockGetDB(users)
	defer dbTest.Close()

	env := controllers.EnvironmentUser{dbTest}

	user, err := env.Db.GetUser(1)
	assert.Equal(t, user, expectedUser1)
	assert.Equal(t, err, nil)

	user, err = env.Db.GetUser(2)
	assert.Equal(t, user, expectedUser2)
	assert.Equal(t, err, nil)

	user, err = env.Db.GetUser(3)
	assert.Equal(t, user, models.User{})
}

func TestLogin(t *testing.T) {
	users := make([]models.User, 0)
	users = append(users, expectedUser1)
	users = append(users, expectedUser2)
	users = append(users, expectedUser3)

	passwords := make([]string, 0)
	for i := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(users[i].Password), 8)
		passwords = append(passwords, string(hashedPassword))
	}
	dbTest, _ := services.NewMockLoginDB(users, passwords)
	defer dbTest.Close()

	env := controllers.EnvironmentUser{dbTest}

	err := env.Db.Login(expectedUser1.Login, expectedUser1.Password)
	assert.Equal(t, err, nil)

	err = env.Db.Login(expectedUser2.Login, expectedUser2.Password)
	assert.Equal(t, err, nil)

	err = env.Db.Login(expectedUser3.Login, expectedUser3.Fullname)
	if err == nil {
		assert.Error(t, err, "Test failed!")
	}
}

func TestDeleteUser(t *testing.T)  {
	users := make([]models.User, 0)
	users = append(users, expectedUser1)
	users = append(users, expectedUser2)
	users = append(users, expectedUser3)

	ids := make([]int, 0)
	ids = append(ids, 1)
	ids = append(ids, 2)
	ids = append(ids, 42)

	dbTest, _ := services.NewMockDeleteDB(users, ids)
	defer dbTest.Close()

	env := controllers.EnvironmentUser{dbTest}


	err := env.Db.DeleteUser(1)
	assert.Equal(t, err, nil)

	err = env.Db.DeleteUser(2)
	assert.Equal(t, err, nil)

	err = env.Db.DeleteUser(24)
	if err == nil {
		assert.Error(t, err, "Test failed!")
	}
}

func TestGetUserHandler(t *testing.T) {

}

