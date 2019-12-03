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

type EnvironmentTask struct {
	Db services.DatastoreTask
}

func(env *EnvironmentTask) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	//проверка токена
	userId, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	//получение информации из url
	var strSlice []string
	idStr := r.URL.Query().Get("id")
	assigneeIdStr := r.URL.Query().Get("assignee_id")
	groupIdStr := r.URL.Query().Get("group_id")
	strSlice = append(strSlice, idStr, assigneeIdStr, groupIdStr)
	title := r.URL.Query().Get("title")
	var idSlice []int

	//запись в массив
	for _, k := range strSlice {
		if k != "" {
			tmp, err := strconv.Atoi(k)
			if err != nil {
				utils.Respond(w, utils.Message(false,"bad parameters", "Bad Request"))
				return
			}
			idSlice = append(idSlice, tmp)
		} else {
			idSlice = append(idSlice, 0)
		}
	}

	//функция возвращает таски по параметрам переданым в url
	tasks, err := env.Db.GetTasks(idSlice, title, userId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	//формирование ответа
	resp := utils.Message(true,"Get tasks", "")
	resp["tasks"]= tasks
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	//проверка токена
	userId, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	//получение task_id
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}
	//получение labels принадлежащих таску
	var labels []models.Label
	task, labels, err := env.Db.GetTaskById(taskId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	//проверка юезра из url и из бд
	if task.CreatorId != userId {
		utils.Respond(w, utils.Message(false, "Id don't match", "Unauthorized"))
		return
	}
	//формирование ответа
	resp := utils.Message(true, "Get task", "")
	resp["task"] = task
	resp["task_labels"] = labels
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)CreateTask(w http.ResponseWriter, r *http.Request) {
	//проверка токена
	id, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	//получение group_id (если нет group_id = 0)
	paramFromURL := mux.Vars(r)
	groupId, _ := strconv.Atoi(paramFromURL["group_id"])
	//получение тела запроса
	task := models.Task{}
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}
	//создание группы
	task.GroupId = groupId
	task, err = env.Db.CreateTask(task, id)
	if err != nil {
		utils.Respond(w, utils.Message(false,err.Error(), "Internal Server Error"))
		return
	}
	//формирование ответа
	resp := utils.Message(true,"Create task", "")
	resp["task"] = task
	utils.Respond(w, resp)
}

//update
func (env *EnvironmentTask)UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	//проверка токена
	userId, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	//получение task_id
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}
	//получение тела запроса
	task := models.Task{}
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}
	//проверка тела запроса
	if  task.Title == "" || task.Description == "" ||
		task.Deadline == 0 || task.State == "" ||
		task.Priority == 0 || task.AssigneeId == 0 || task.Duration == 0 {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}
	//обновление
	task, err = env.Db.UpdateTask(task, taskId, userId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(),"Internal Server Error"))
		return
	}
	//формирование ответа
	resp := utils.Message(true, "Update task", "")
	resp["task"] = task
	utils.Respond(w, resp)
}
