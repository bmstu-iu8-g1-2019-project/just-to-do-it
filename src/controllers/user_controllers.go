package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/auth"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/utils"
)

type EnvironmentUser struct {
	Db services.DatastoreUser
}

// handle for user authorization
func (env *EnvironmentUser) ResponseLoginHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	// write from the received data to the User structure
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	// function checks login and password in database
	user, err = env.Db.Login(user.Login, user.Password)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid login or password","Unauthorized"))
		return
	}

	accToken, err := auth.CreateAccessToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForAccToken(w, accToken)
	refToken, err := auth.CreateRefreshToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForRefToken(w, refToken)

	resp := utils.Message(true, "Logged In", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) ResponseRegisterHandler (w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	if user.Password == ""  || user.Login == "" || user.Email == ""{
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	//function that detects the user and hashes his password
	user, msg, errStr := env.Db.Register(user)
	if msg != "" {
		utils.Respond(w, utils.Message(false, msg, errStr))
		return
	}

	accToken, err := auth.CreateAccessToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForAccToken(w, accToken)
	refToken, err := auth.CreateRefreshToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForRefToken(w, refToken)

	resp := utils.Message(true, "User created", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) ConfirmEmailHandler (w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	err := env.Db.Confirm(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) GetUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	//проверка и в случае таймута рефреш токена
	//resp := auth.CheckTokenAndRefresh(w, r, id)
	//if resp["status"] == false {
	//	utils.Respond(w, resp)
	//	return
	//}
	err = auth.TokenValid(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	//Get user
	user, err := env.Db.GetUser(int(id))
	if err != nil {
		utils.Respond(w, utils.Message(false,"Not found user in db","Internal Server Error"))
		return
	}
	resp := utils.Message(true, "Get user", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) UpdateUserHandler (w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	//получаем из запроса структуру
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	if  user.Email == "" || user.Fullname == "" ||
		user.Login == "" || user.Password == "" {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}
	//получаем из url id
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	//проверка и в случае таймута рефреш токена
	//resp := auth.CheckTokenAndRefresh(w, r, id)
	//if resp["status"] == false {
	//	utils.Respond(w, resp)
	//	return
	//}

	// обновляем юзера по пришедшему из path id и json
	user, err = env.Db.UpdateUser(int(id), user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Database error", "Internal Server Error"))
		return
	}
	resp := utils.Message(true, "Update user", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}
	//проверка и в случае таймута рефреш токена
	//resp := auth.CheckTokenAndRefresh(w, r, id)
	//if resp["status"] == false {
	//	utils.Respond(w, resp)
	//	return
	//}

	err = env.Db.DeleteUser(id)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Database error","Internal Server Error"))
		return
	}
	resp := utils.Message(true, "User deleted", "")
	utils.Respond(w, resp)
}
