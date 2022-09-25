package services

import (
	"github.com/JECSand/go-rest-api-boilerplate/models"
)

// UserService is an interface used to manage the relevant user doc controllers
type UserService interface {
	AuthenticateUser(u *models.User) (*models.User, error)
	UpdatePassword(u *models.User, CurrentPassword string, newPassword string) (*models.User, error)
	UserCreate(u *models.User) (*models.User, error)
	UserDelete(u *models.User) (*models.User, error)
	UsersFind(u *models.User) ([]*models.User, error)
	UserFind(u *models.User) (*models.User, error)
	UserUpdate(u *models.User) (*models.User, error)
	UserDocInsert(u *models.User) (*models.User, error)
}
