package server

import (
	"encoding/json"
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/auth"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"github.com/JECSand/go-rest-api-boilerplate/services"
	"github.com/JECSand/go-rest-api-boilerplate/utilities"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type groupRouter struct {
	aService *services.TokenService
	gService services.GroupService
	uService services.UserService
	tService services.TaskService
	fService services.FileService
}

// NewGroupRouter is a function that initializes a new groupRouter struct
func NewGroupRouter(router *mux.Router, a *services.TokenService, g services.GroupService, u services.UserService, t services.TaskService, f services.FileService) *mux.Router {
	gRouter := groupRouter{a, g, u, t, f}
	router.HandleFunc("/groups", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/groups", a.AdminTokenVerifyMiddleWare(gRouter.GroupsShow)).Methods("GET")
	router.HandleFunc("/groups", a.RootAdminTokenVerifyMiddleWare(gRouter.CreateGroup)).Methods("POST")
	router.HandleFunc("/groups/{groupId}", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/groups/{groupId}", a.AdminTokenVerifyMiddleWare(gRouter.GroupShow)).Methods("GET")
	router.HandleFunc("/groups/{groupId}", a.RootAdminTokenVerifyMiddleWare(gRouter.DeleteGroup)).Methods("DELETE")
	router.HandleFunc("/groups/{groupId}", a.AdminTokenVerifyMiddleWare(gRouter.ModifyGroup)).Methods("PATCH")
	router.HandleFunc("/groups/{groupId}/users", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/groups/{groupId}/users", a.AdminTokenVerifyMiddleWare(gRouter.GroupUsersShow)).Methods("GET")
	return router
}

// getGroupUsers asynchronously gets a group and its users from the database
func (gr *groupRouter) getGroupUsers(groupId string) (*groupUsersDTO, error) {
	var dto groupUsersDTO
	gOutCh := make(chan *models.Group)
	gErrCh := make(chan error)
	uOutCh := make(chan []*models.User)
	uErrCh := make(chan error)
	go func() {
		reG, err := gr.gService.GroupFind(&models.Group{Id: groupId})
		gOutCh <- reG
		gErrCh <- err
	}()
	go func() {
		reU, err := gr.uService.UsersFind(&models.User{GroupId: groupId})
		uOutCh <- reU
		uErrCh <- err
	}()
	for i := 0; i < 4; i++ {
		select {
		case gOut := <-gOutCh:
			dto.Group = gOut
		case gErr := <-gErrCh:
			if gErr != nil {
				return &dto, gErr
			}
		case uOut := <-uOutCh:
			dto.Users = uOut
		case uErr := <-uErrCh:
			if uErr != nil {
				return &dto, uErr
			}
		}
	}
	return &dto, nil
}

// GroupUsersShow returns a groupUsersDTO
func (gr *groupRouter) GroupUsersShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	groupId := vars["groupId"]
	if !utilities.CheckObjectID(groupId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing groupId"})
		return
	}
	groupId, err = auth.VerifyGroupRequestScope(r, groupId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	dto, err := gr.getGroupUsers(groupId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	dto.clean()
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		return
	}
	return
}

// GroupsShow returns all groups to client
func (gr *groupRouter) GroupsShow(w http.ResponseWriter, r *http.Request) {
	w = utilities.SetResponseHeaders(w, "", "")
	tokenData, err := auth.LoadTokenFromRequest(r)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	groups, err := gr.gService.GroupsFind(tokenData.GetGroupsScope())
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&groupsDTO{Groups: groups}); err != nil {
		return
	}
}

// CreateGroup from a REST Request post body
func (gr *groupRouter) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &group); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	group.Id = utilities.GenerateObjectID()
	group.RootAdmin = false
	g, err := gr.gService.GroupCreate(&group)
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

// ModifyGroup to update a group document
func (gr *groupRouter) ModifyGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	if !utilities.CheckObjectID(groupId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing groupId"})
		return
	}
	var group models.Group
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &group); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	groupId, err = auth.VerifyGroupRequestScope(r, groupId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	group.Id = groupId
	g, err := gr.gService.GroupUpdate(&group)
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

// GroupShow shows a specific group
func (gr *groupRouter) GroupShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	groupId := vars["groupId"]
	if !utilities.CheckObjectID(groupId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing groupId"})
		return
	}
	groupId, err = auth.VerifyGroupRequestScope(r, groupId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	group, err := gr.gService.GroupFind(&models.Group{Id: groupId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(group); err != nil {
		return
	}
	return
}

// DeleteGroup deletes a group
func (gr *groupRouter) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	if !utilities.CheckObjectID(groupId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing groupId"})
		return
	}
	groupUsers, err := gr.getGroupUsers(groupId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	err = gr.deleteGroupAssets(groupUsers.Group, groupUsers.Users)
	if err != nil {
		utilities.RespondWithError(w, http.StatusInternalServerError, utilities.JWTError{Message: err.Error()})
		return
	}
	group, err := gr.gService.GroupDelete(&models.Group{Id: groupId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(group); err != nil {
		return
	}
	return
}

// deleteGroupAssets asynchronously gets a group and its users from the database
func (gr *groupRouter) deleteGroupAssets(group *models.Group, users []*models.User) error {
	if !group.CheckID("id") {
		return errors.New("filter id cannot be empty for mass delete")
	}
	fErrCh := make(chan error) // Images Files Bulk Delete
	uErrCh := make(chan error) // Delete Group Users
	tErrCh := make(chan error) // Delete Group Tasks
	go func() {
		err := gr.fService.FileDeleteMany(models.UsersToFiles(users))
		fErrCh <- err
	}()
	go func() {
		_, err := gr.uService.UserDeleteMany(&models.User{GroupId: group.Id})
		uErrCh <- err
	}()
	go func() {
		_, err := gr.tService.TaskDeleteMany(&models.Task{GroupId: group.Id})
		tErrCh <- err
	}()
	for i := 0; i < 3; i++ {
		select {
		case fErr := <-fErrCh:
			if fErr != nil {
				return fErr
			}
		case uErr := <-uErrCh:
			if uErr != nil {
				return uErr
			}
		case tErr := <-tErrCh:
			if tErr != nil {
				return tErr
			}
		}
	}
	return nil
}
