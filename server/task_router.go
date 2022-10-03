package server

import (
	"encoding/json"
	"github.com/JECSand/go-rest-api-boilerplate/auth"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"github.com/JECSand/go-rest-api-boilerplate/services"
	"github.com/JECSand/go-rest-api-boilerplate/utilities"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type taskRouter struct {
	aService *services.TokenService
	tService services.TaskService
}

// NewTaskRouter is a function that initializes a new groupRouter struct
func NewTaskRouter(router *mux.Router, a *services.TokenService, t services.TaskService) *mux.Router {
	gRouter := taskRouter{a, t}
	router.HandleFunc("/tasks", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/tasks", a.MemberTokenVerifyMiddleWare(gRouter.TasksShow)).Methods("GET")
	router.HandleFunc("/tasks", a.MemberTokenVerifyMiddleWare(gRouter.CreateTask)).Methods("POST")
	router.HandleFunc("/tasks/{taskId}", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/tasks/{taskId}", a.MemberTokenVerifyMiddleWare(gRouter.TaskShow)).Methods("GET")
	router.HandleFunc("/tasks/{taskId}", a.MemberTokenVerifyMiddleWare(gRouter.DeleteTask)).Methods("DELETE")
	router.HandleFunc("/tasks/{taskId}", a.MemberTokenVerifyMiddleWare(gRouter.ModifyTask)).Methods("PATCH")
	return router
}

// TasksShow returns all tasks to client
func (gr *taskRouter) TasksShow(w http.ResponseWriter, r *http.Request) {
	var filter models.Task
	userScope, err := auth.VerifyRequestScope(r)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	filter.LoadScope(userScope)
	tasks, err := gr.tService.TasksFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&tasksDTO{Tasks: tasks}); err != nil {
		return
	}
}

// CreateTask from a REST Request post body
func (gr *taskRouter) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &task); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	userScope, err := auth.VerifyRequestScope(r)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	task.LoadScope(userScope)
	task.Id = utilities.GenerateObjectID()
	if !task.CheckID("user_id") || !task.CheckID("group_id") {
		td, err := auth.LoadTokenFromRequest(r)
		if err != nil {
			utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
			return
		}
		if !task.CheckID("user_id") {
			task.UserId = td.UserId
		}
		if !task.CheckID("group_id") {
			task.GroupId = td.GroupId
		}
	}
	g, err := gr.tService.TaskCreate(&task)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	} else {
		w = utilities.SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(g); err != nil {
			return
		}
	}
}

// ModifyTask to update a task document
func (gr *taskRouter) ModifyTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["taskId"]
	if !utilities.CheckObjectID(taskId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing taskId"})
		return
	}
	var task models.Task
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &task); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	task.Id = taskId
	g, err := gr.tService.TaskUpdate(&task)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	} else {
		w = utilities.SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		if err = json.NewEncoder(w).Encode(g); err != nil {
			return
		}
	}
}

// TaskShow shows a specific task
func (gr *taskRouter) TaskShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["taskId"]
	if !utilities.CheckObjectID(taskId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing taskId"})
		return
	}
	var filter models.Task
	userScope, err := auth.VerifyRequestScope(r)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	filter.LoadScope(userScope)
	filter.Id = taskId
	task, err := gr.tService.TaskFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(task); err != nil {
		return
	}
	return
}

// DeleteTask deletes a task
func (gr *taskRouter) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["taskId"]
	if !utilities.CheckObjectID(taskId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing taskId"})
		return
	}
	var filter models.Task
	userScope, err := auth.VerifyRequestScope(r)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	filter.LoadScope(userScope)
	filter.Id = taskId
	task, err := gr.tService.TaskDelete(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(task); err != nil {
		return
	}
	return
}
