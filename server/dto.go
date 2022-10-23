package server

import (
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/models"
)

/*
================ User DTOs ==================
*/

// updatePassword is used when updating a user password
type updatePassword struct {
	NewPassword     string `json:"new_password"`
	CurrentPassword string `json:"current_password"`
}

// userSignIn is used when updating a user password
type userSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// toUser converts userSignIn DTO to a user
func (u *userSignIn) toUser() (*models.User, error) {
	if u.Email == "" {
		return &models.User{}, errors.New("missing user email")
	}
	if u.Password == "" {
		return &models.User{}, errors.New("missing user password")
	}
	return &models.User{
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

// usersDTO is used when returning a slice of User
type usersDTO struct {
	Users []*models.User `json:"users"`
}

// clean ensures the users in the usersDTO have no passwords set
func (u *usersDTO) clean() {
	for i, _ := range u.Users {
		u.Users[i].Password = ""
	}
}

// userTasksDTO is used when returning user with associated tasks
type userTasksDTO struct {
	User  *models.User   `json:"user"`
	Tasks []*models.Task `json:"tasks"`
}

// clean ensures the users in the userTasksDTO have password set
func (u *userTasksDTO) clean() {
	u.User.Password = ""
}

/*
================ Group DTOs ==================
*/

// groupsDTO is used when returning a slice of Group
type groupsDTO struct {
	Groups []*models.Group `json:"groups"`
}

// groupUsersDTO is used when returning a group with its associated users
type groupUsersDTO struct {
	Group *models.Group  `json:"group"`
	Users []*models.User `json:"users"`
}

// clean ensures the users in the groupUsersDTO have no passwords set
func (u *groupUsersDTO) clean() {
	for i, _ := range u.Users {
		u.Users[i].Password = ""
	}
}

// groupTasksDTO is used when returning a group with its associated tasks
type groupTasksDTO struct {
	Group *models.Group  `json:"group"`
	Tasks []*models.Task `json:"tasks"`
}

/*
================ Task DTOs ==================
*/

// tasksDTO is used when returning a slice of Task
type tasksDTO struct {
	Tasks []*models.Task `json:"tasks"`
}
