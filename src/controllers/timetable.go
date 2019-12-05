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

type EnvironmentTimeTable struct {
	Db services.DatastoreTimeTable
}

func (env *EnvironmentTimeTable)GetTimetableHandler(w http.ResponseWriter, r *http.Request) {
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

	tables, err := env.Db.GetTimetables(params)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Got timetables!", "")
	resp["timetables"] = tables
	utils.Respond(w, resp)
}

func (env *EnvironmentTimeTable)UpdateTimetableHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid id", "Bad Request"))
		return
	}

	table := models.Timetable{}
	err = json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid body", "Bad Request"))
		return
	}

	if table.Id <= 0 || table.GroupId <= 0 || table.BeginInterval <= 0 || table.EndInterval <= 0 {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	table, err = env.Db.UpdateTimetable(id, table)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Database error", "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update timetable", "")
	resp["timetable"] = table
	utils.Respond(w, resp)
}

func (env *EnvironmentTimeTable) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	err = env.Db.DeleteTimetable(id)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Database error","Internal Server Error"))
		return
	}
	resp := utils.Message(true, "Timetable deleted", "")
	utils.Respond(w, resp)
}
