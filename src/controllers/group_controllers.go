package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/utils"
)

type EnvironmentGroup struct {
	Db services.DatastoreGroup
}

func (env *EnvironmentGroup)CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	//получение тела из запроса
	group := models.Group{}
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}
	if group.Title == "" || group.Description == "" {

	}
	//создание группы
	group, err = env.Db.CreateGroup(group)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	//формирование ответа
	resp := utils.Message(true, "Created group", "")
	resp["group"] = group
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	//получение group_id
	paramFromURL := mux.Vars(r)
	groupId, err := strconv.Atoi(paramFromURL["group_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}
	//получение информации о группе
	group, err := env.Db.GetGroup(groupId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	//формирование ответа
	resp := utils.Message(true, "Get group", "")
	resp["group"] = group
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	//получение group_id
	paramFromURL := mux.Vars(r)
	groupId, err := strconv.Atoi(paramFromURL["group_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}
	//получение тела из запроса
	group := models.Group{}
	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid id", "Bad Request"))
		return
	}
	if group.Title == "" || group.Description == "" {
		utils.Respond(w, utils.Message(false, "Invalid body", "Bad Request"))
		return
	}
	//обновление
	group, err = env.Db.UpdateGroup(groupId, group)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	//формирование ответа
	resp := utils.Message(true, "Update group", "")
	group.Id = groupId
	resp["group"] = group
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	//получение group_id
	paramFromURL := mux.Vars(r)
	groupId, err := strconv.Atoi(paramFromURL["group_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}
	//удаление группы и всех задач связанных с группой
	err = env.Db.DeleteGroup(groupId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	//формирование ответа
	resp := utils.Message(true, "Deleted group", "")
	utils.Respond(w, resp)
}
