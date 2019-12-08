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

func (env *EnvironmentScope)CreateScopeHandler(w http.ResponseWriter, r* http.Request) {
	// Получение user_id и group_id
	paramsFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramsFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Bad parameters", "Bad Request"))
		return
	}
	groupId, err := strconv.Atoi(paramsFromURL["group_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Bad parameters", "Bad Request"))
		return
	}
	// Получение тела запроса и проверка
	scope := models.Scope{}
	err = json.NewDecoder(r.Body).Decode(&scope)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}
	scope.CreatorId = id
	scope.GroupId = groupId
	err = models.ValidTimetable(scope)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}
	// Запись в бд
	scope, err = env.Db.CreateScope(scope)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	// Формирование ответа
	resp := utils.Message(true, "Created scope", "")
	resp["scope"] = scope
	utils.Respond(w, resp)
}

func (env *EnvironmentScope)AddTaskInScopeHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	scopeId, err := strconv.Atoi(paramFromURL["scope_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Bad parameters", "Bad Request"))
		return
	}
	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Bad parameters", "Bad Request"))
		return
	}

	timetable, err := env.Db.AddTaskInScope(scopeId, taskId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Added task", "")
	resp["timetable"] = timetable
	utils.Respond(w, resp)
}

func (env *EnvironmentScope)GetScopesWithIntervalHandler(w http.ResponseWriter, r *http.Request) {
	// Получение и проверка границ
	beginInterval := r.URL.Query().Get("begin")
	endInterval := r.URL.Query().Get("end")
	if beginInterval == "" && endInterval == "" {
		utils.Respond(w, utils.Message(false, "Bad parameters", "Bad Request"))
		return
	}
	begin, err := strconv.Atoi(beginInterval)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Bad parameters", "Bad Request"))
		return
	}
	end, err := strconv.Atoi(endInterval)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Bad parameters", "Bad Request"))
		return
	}
	// Получение scope'ов в интервале от begin до end
	scopes, err := env.Db.GetScopesWithInterval(begin, end)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	// Получение тасков принадлежащих scope'ам
	tasks := make([]models.Task, 0)
	for _, scope := range scopes {
		tasksInScope, err := env.Db.GetTasksFromScope(scope.Id)
		if err != nil {
			utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
			return
		}
		for _, i := range tasksInScope {
			tasks = append(tasks, i)
		}
	}
	// Формирование ответа
	resp := utils.Message(true, "Get scopes", "")
	resp["scopes"] = scopes
	resp["tasks"] = tasks
	utils.Respond(w, resp)
}

func (env *EnvironmentScope)GetScopesHandler(w http.ResponseWriter, r *http.Request) {
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
	id, err := strconv.Atoi(paramFromURL["scope_id"])
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

	if  scope.GroupId <= 0 || scope.BeginInterval <= 0 || scope.EndInterval <= 0 {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	scope, err = env.Db.UpdateScope(id, scope)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update scope", "")
	resp["scope"] = scope
	utils.Respond(w, resp)
}

func (env *EnvironmentScope)DeleteScopeHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["scope_id"])
	err = env.Db.DeleteScope(id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(),"Internal Server Error"))
		return
	}
	resp := utils.Message(true, "Scope deleted", "")
	utils.Respond(w, resp)
}

func (env *EnvironmentScope)CreateSmartScopeHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["scope_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Bad parameters", "Bad Request"))
		return
	}

	tasks, err := env.Db.CreateSmartScope(id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(),"Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Created scope", "")
	resp["Tasks"] = tasks
	utils.Respond(w, resp)
}
