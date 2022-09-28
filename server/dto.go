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

/*
================ Group DTOs ==================
*/

// groupsDTO is used when returning a slice of Group
type groupsDTO struct {
	Groups []*models.Group `json:"groups"`
}

/*
================ Task DTOs ==================
*/

// tasksDTO is used when returning a slice of Task
type tasksDTO struct {
	Tasks []*models.Task `json:"tasks"`
}
