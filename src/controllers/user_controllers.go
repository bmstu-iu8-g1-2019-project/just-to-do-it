package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/auth"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/utils"
)

type EnvironmentUser struct {
	Db services.DatastoreUser
}

func (env *EnvironmentUser) ResponseLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Получение логина и пароля из тела запроса
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	// Функция проверяет логин и пароль в бд и возвращает информацию о пользователе
	user, err = env.Db.Login(user.Login, user.Password)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid login or password","Unauthorized"))
		return
	}
	// Генерация access токена
	accToken, err := auth.CreateAccessToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForAccToken(w, accToken)
	// Генарция refresh токена
	refToken, err := auth.CreateRefreshToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForRefToken(w, refToken)
	// Формирование ответа
	resp := utils.Message(true, "Logged In", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) ResponseRegisterHandler (w http.ResponseWriter, r *http.Request) {
	// получение json'a
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	// проверка полей
	if user.Password == ""  || user.Login == "" || user.Email == "" {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	// функция добавляет пользователя в бд и хэширует его пароль
	user, msg, errStr := env.Db.Register(user)
	if msg != "" {
		utils.Respond(w, utils.Message(false, msg, errStr))
		return
	}
	// Генераци access токена
	accToken, err := auth.CreateAccessToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForAccToken(w, accToken)
	// Генерация refresh токена
	refToken, err := auth.CreateRefreshToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForRefToken(w, refToken)
	// Формирование ответа
	resp := utils.Message(true, "User created", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) ConfirmEmailHandler (w http.ResponseWriter, r *http.Request) {
	// получние хэша из ссылки
	hash := r.URL.Query().Get("hash")
	// функция в случае успешного подтверждения изменяет
	// поле acc_verified у юзера на true
	// в случае перехода по старой ссылке отправляет новое письмо
	err := env.Db.Confirm(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) GetUserHandler (w http.ResponseWriter, r *http.Request) {
	//проверка токена
	id, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	// функция возвращает из бд информацию о юезере
	user, err := env.Db.GetUser(int(id))
	if err != nil {
		utils.Respond(w, utils.Message(false,"Not found user in db","Internal Server Error"))
		return
	}
	// формирование ответа
	resp := utils.Message(true, "Get user", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) UpdateUserHandler (w http.ResponseWriter, r *http.Request) {
	//проверка токена
	id, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	// получаем из запроса структуру
	user := models.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	// проверка тела json'a
	if  user.Email == "" || user.Fullname == "" ||
		user.Login == "" || user.Password == "" {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	// обновляем юзера по пришедшему из path id и json
	user, err = env.Db.UpdateUser(int(id), user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Database error", "Internal Server Error"))
		return
	}
	// формирование ответа
	resp := utils.Message(true, "Update user", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	//проверка токена
	id, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	// удаление из бд пользователя
	err = env.Db.DeleteUser(id)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Database error","Internal Server Error"))
		return
	}
	// формирование ответа
	resp := utils.Message(true, "User deleted", "")
	utils.Respond(w, resp)
}