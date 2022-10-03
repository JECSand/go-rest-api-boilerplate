package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/auth"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"github.com/JECSand/go-rest-api-boilerplate/services"
	"github.com/JECSand/go-rest-api-boilerplate/utilities"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"os"
	"time"
)

type userRouter struct {
	aService *services.TokenService
	uService services.UserService
	gService services.GroupService
	tService services.TaskService
	fService services.FileService
}

// NewUserRouter is a function that initializes a new userRouter struct
func NewUserRouter(router *mux.Router, a *services.TokenService, u services.UserService, g services.GroupService, t services.TaskService, f services.FileService) *mux.Router {
	uRouter := userRouter{a, u, g, t, f}
	router.HandleFunc("/auth", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth", uRouter.SignIn).Methods("POST")
	router.HandleFunc("/auth", a.MemberTokenVerifyMiddleWare(uRouter.RefreshSession)).Methods("GET")
	router.HandleFunc("/auth", a.MemberTokenVerifyMiddleWare(uRouter.SignOut)).Methods("DELETE")
	router.HandleFunc("/auth/register", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/register", uRouter.RegisterUser).Methods("POST")
	router.HandleFunc("/auth/api-key", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/api-key", a.MemberTokenVerifyMiddleWare(uRouter.GenerateAPIKey)).Methods("GET")
	router.HandleFunc("/auth/password", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/password", a.MemberTokenVerifyMiddleWare(uRouter.UpdatePassword)).Methods("POST")
	router.HandleFunc("/users", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users", a.MemberTokenVerifyMiddleWare(uRouter.UsersShow)).Methods("GET")
	router.HandleFunc("/users/{userId}", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users/{userId}", a.MemberTokenVerifyMiddleWare(uRouter.UserShow)).Methods("GET")
	router.HandleFunc("/users", a.AdminTokenVerifyMiddleWare(uRouter.CreateUser)).Methods("POST")
	router.HandleFunc("/users/{userId}", a.AdminTokenVerifyMiddleWare(uRouter.DeleteUser)).Methods("DELETE")
	router.HandleFunc("/users/{userId}", a.MemberTokenVerifyMiddleWare(uRouter.ModifyUser)).Methods("PATCH")
	router.HandleFunc("/users/{userId}/image", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users/{userId}/image", a.MemberTokenVerifyMiddleWare(uRouter.UploadImage)).Methods("POST")
	router.HandleFunc("/users/{userId}/image", a.MemberTokenVerifyMiddleWare(uRouter.GetImage)).Methods("GET")
	return router
}

// UpdatePassword is the handler function that manages the user password update process
func (ur *userRouter) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	decodedToken, err := auth.DecodeJWT(r.Header.Get("Auth-Token"))
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	var pw updatePassword
	err = json.Unmarshal(body, &pw)
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	inUser := decodedToken.ToUser()
	u, err := ur.uService.UpdatePassword(inUser, pw.CurrentPassword, pw.NewPassword)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	} else {
		w = utilities.SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		u.Password = ""
		if err = json.NewEncoder(w).Encode(u); err != nil {
			return
		}
		return
	}
}

// ModifyUser is the handler function that updates a user
func (ur *userRouter) ModifyUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if !utilities.CheckObjectID(userId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing userId"})
		return
	}
	var user models.User
	user.Id = userId
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &user); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	userScope, err := auth.VerifyUserRequestScope(r, userId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	user.LoadScope(userScope, "update")
	u, err := ur.uService.UserUpdate(&user)
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	} else {
		w = utilities.SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		if err = json.NewEncoder(w).Encode(u); err != nil {
			return
		}
		return
	}
}

// SignIn is the handler function that manages the user SignIn process
func (ur *userRouter) SignIn(w http.ResponseWriter, r *http.Request) {
	var dto userSignIn
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &dto); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	user, err := dto.toUser()
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	u, err := ur.uService.AuthenticateUser(user)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	} else {
		sessionToken, err := ur.aService.GenerateToken(u, "session")
		if err != nil {
			utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
			return
		}
		w = utilities.SetResponseHeaders(w, sessionToken, "")
		w.WriteHeader(http.StatusOK)
		u.Password = ""
		if err = json.NewEncoder(w).Encode(u); err != nil {
			return
		}
		return
	}
}

// RefreshSession is the handler function that refreshes a users JWT token
func (ur *userRouter) RefreshSession(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData, err := auth.DecodeJWT(authToken)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	user, err := ur.uService.UserFind(tokenData.ToUser())
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	newToken, err := ur.aService.GenerateToken(user, "session")
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, newToken, "")
	w.WriteHeader(http.StatusOK)
	return
}

// GenerateAPIKey is the handler function that generates 6 month API Key for a given user
func (ur *userRouter) GenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData, err := auth.DecodeJWT(authToken)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	user, err := ur.uService.UserFind(tokenData.ToUser())
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	apiKey, err := ur.aService.GenerateToken(user, "api")
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", apiKey)
	w.WriteHeader(http.StatusOK)
	return
}

// SignOut is the handler function that ends a users session
func (ur *userRouter) SignOut(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	err := ur.aService.BlacklistAuthToken(authToken)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	return
}

// RegisterUser handler function that registers a new user
func (ur *userRouter) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("REGISTRATION") == "OFF" {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: "Not Found"})
		return
	} else {
		var user models.User
		body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
			return
		}
		if err = r.Body.Close(); err != nil {
			utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
			return
		}
		if err = json.Unmarshal(body, &user); err != nil {
			utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
			return
		}
		var group models.Group
		groupName := user.Email
		groupName += "_group"
		group.Name = groupName
		group.Id = utilities.GenerateObjectID()
		group.RootAdmin = false
		g, err := ur.gService.GroupCreate(&group)
		if err != nil {
			utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
			return
		}
		user.Role = "admin"
		user.GroupId = g.Id
		u, err := ur.uService.UserCreate(&user)
		if err != nil {
			utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
			return
		} else {
			newToken, err := ur.aService.GenerateToken(u, "session")
			if err != nil {
				utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
				return
			}
			w = utilities.SetResponseHeaders(w, newToken, "")
			w.WriteHeader(http.StatusCreated)
			u.Password = ""
			if err = json.NewEncoder(w).Encode(u); err != nil {
				return
			}
			return
		}
	}
}

// CreateUser is the handler function that creates a new user
func (ur *userRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &user); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	decodedToken, err := auth.DecodeJWT(r.Header.Get("Auth-Token"))
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	userScope := decodedToken.GetUsersScope()
	user.LoadScope(userScope, "create")
	u, err := ur.uService.UserCreate(&user)
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	} else {
		w = utilities.SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusCreated)
		u.Password = ""
		if err = json.NewEncoder(w).Encode(u); err != nil {
			return
		}
		return
	}
}

// UsersShow is the handler that shows a specific user
func (ur *userRouter) UsersShow(w http.ResponseWriter, r *http.Request) {
	decodedToken, err := auth.DecodeJWT(r.Header.Get("Auth-Token"))
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	var filter models.User
	userScope := decodedToken.GetUsersScope()
	filter.LoadScope(userScope, "find")
	users, err := ur.uService.UsersFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	dto := usersDTO{Users: users}
	dto.clean()
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		return
	}
	return
}

// UserShow is the handler that shows all users
func (ur *userRouter) UserShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if !utilities.CheckObjectID(userId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing userId"})
		return
	}
	filter := models.User{Id: userId}
	userScope, err := auth.VerifyUserRequestScope(r, userId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	filter.LoadScope(userScope, "find")
	user, err := ur.uService.UserFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	user.Password = ""
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		return
	}
	return
}

// DeleteUser is the handler function that deletes a user
func (ur *userRouter) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if !utilities.CheckObjectID(userId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing userId"})
		return
	}
	filter := models.User{Id: userId}
	userScope, err := auth.VerifyUserRequestScope(r, userId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	filter.LoadScope(userScope, "find")
	user, err := ur.uService.UserFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	err = ur.deleteUserAssets(user)
	if err != nil {
		utilities.RespondWithError(w, http.StatusInternalServerError, utilities.JWTError{Message: err.Error()})
		return
	}
	user, err = ur.uService.UserDelete(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	if user.Id != "" {
		w = utilities.SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode("User Deleted"); err != nil {
			return
		}
		return
	}
}

// UploadImage allows for a user image to be associated with the User record
func (ur *userRouter) UploadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if !utilities.CheckObjectID(userId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing userId"})
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	defer file.Close()
	filter := models.User{Id: userId}
	userScope, err := auth.VerifyUserRequestScope(r, userId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	filter.LoadScope(userScope, "find")
	user, err := ur.uService.UserFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: "user not found"})
		return
	}
	newImage := false
	if !user.CheckID("image_id") {
		user.ImageId = utilities.GenerateObjectID()
		newImage = true
	}
	f := &models.File{Id: user.ImageId, OwnerType: "user", OwnerId: user.Id, BucketType: "user-images", Name: handler.Filename}
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		utilities.RespondWithError(w, http.StatusInternalServerError, utilities.JWTError{Message: err.Error()})
		return
	}
	if newImage {
		f, err = ur.fService.FileCreate(f, buf.Bytes())
		if err != nil {
			utilities.RespondWithError(w, http.StatusInternalServerError, utilities.JWTError{Message: err.Error()})
			return
		}
		user, err = ur.uService.UserUpdate(&models.User{Id: user.Id, ImageId: user.ImageId})
		if err != nil {
			utilities.RespondWithError(w, http.StatusInternalServerError, utilities.JWTError{Message: err.Error()})
			return
		}
	} else {
		f, err = ur.fService.FileUpdate(f, buf.Bytes())
		if err != nil {
			utilities.RespondWithError(w, http.StatusInternalServerError, utilities.JWTError{Message: err.Error()})
			return
		}
	}
	user.Password = ""
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		return
	}
	return
}

// GetImage returns the file contents of a User image
func (ur *userRouter) GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if !utilities.CheckObjectID(userId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing userId"})
		return
	}
	filter := models.User{Id: userId}
	userScope, err := auth.VerifyUserRequestScope(r, userId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	filter.LoadScope(userScope, "find")
	user, err := ur.uService.UserFind(&filter)
	if err != nil || !user.CheckID("image_id") {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: "user image not found"})
		return
	}
	file, err := ur.fService.FileFind(&models.File{Id: user.ImageId})
	if err != nil || file.OwnerType != "user" || file.OwnerId != user.Id {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: "unauthorized"})
		return
	}
	contents, err := ur.fService.RetrieveFile(&models.File{GridFSId: file.GridFSId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	modTime := time.Now()
	cd := mime.FormatMediaType("attachment", map[string]string{"filename": file.Name})
	w.Header().Set("Content-Disposition", cd)
	w.Header().Set("Content-Type", "application/octet-stream")
	contentReader := bytes.NewReader(contents.Bytes())
	http.ServeContent(w, r, file.Name, modTime, contentReader)
}

// deleteUserAssets asynchronously gets a group and its users from the database
func (ur *userRouter) deleteUserAssets(user *models.User) error {
	if !user.CheckID("id") {
		return errors.New("filter id cannot be empty for mass delete")
	}
	gErrCh := make(chan error)
	uErrCh := make(chan error)
	go func() {
		if user.CheckID("image_id") {
			_, err := ur.fService.FileDelete(&models.File{OwnerId: user.Id, OwnerType: "user"})
			gErrCh <- err
		} else {
			gErrCh <- nil
		}
	}()
	go func() {
		_, err := ur.tService.TaskDeleteMany(&models.Task{UserId: user.Id})
		uErrCh <- err
	}()
	for i := 0; i < 2; i++ {
		select {
		case gErr := <-gErrCh:
			if gErr != nil {
				return gErr
			}
		case uErr := <-uErrCh:
			if uErr != nil {
				return uErr
			}
		}
	}
	return nil
}
