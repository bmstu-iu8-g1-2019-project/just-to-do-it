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

type EnvironmentScope struct {
	Db services.DatastoreScope
}

func (env *EnvironmentScope)CreateScope(w http.ResponseWriter, r* http.Request) {
	paramsFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramsFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Bad parameters", "Bad Request"))
		return
	}

	scope := models.Scope{}
	err = json.NewDecoder(r.Body).Decode(&scope)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}

	scope, err = env.Db.CreateScope(id, scope)
	if err != nil {
		utils.Respond(w, utils.Message(false,err.Error(), "Internal Server Error"))
		return
	}
}

func (env *EnvironmentScope)GetScopeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	creatorId := r.URL.Query().Get("creator_id")
	groupId := r.URL.Query().Get("group_id")

	suspects := make([]string, 0)
	suspects = append(suspects, id, creatorId, groupId)

	params := make([]int, 0)
	for _, value := range suspects {
		if value != "" {
			tmp, err := strconv.Atoi(value)
			if err != nil {
				utils.Respond(w, utils.Message(false,"Bad parameters", "Bad Request"))
				return
			}
			params = append(params, tmp)
		} else {
			params = append(params, 0)
		}
	}

	scopes, err := env.Db.GetScopes(params)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Got scopes!", "")
	resp["scopes"] = scopes
	utils.Respond(w, resp)
}

func (env *EnvironmentScope)UpdateScopeHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid id", "Bad Request"))
		return
	}

	scope := models.Scope{}
	err = json.NewDecoder(r.Body).Decode(&scope)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid body", "Bad Request"))
		return
	}

	if scope.Id <= 0 || scope.GroupId <= 0 || scope.BeginInterval <= 0 || scope.EndInterval <= 0 {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	scope, err = env.Db.UpdateScope(id, scope)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Database error", "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update scope", "")
	resp["scope"] = scope
	utils.Respond(w, resp)
}

func (env *EnvironmentScope) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	err = env.Db.DeleteScope(id)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Database error","Internal Server Error"))
		return
	}
	resp := utils.Message(true, "Scope deleted", "")
	utils.Respond(w, resp)
}
