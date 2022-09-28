package services

import "github.com/JECSand/go-rest-api-boilerplate/models"

// GroupService is an interface used to manage the relevant group doc controllers
type GroupService interface {
	GroupCreate(g *models.Group) (*models.Group, error)
	GroupFind(g *models.Group) (*models.Group, error)
	GroupsFind(g *models.Group) ([]*models.Group, error)
	GroupDelete(g *models.Group) (*models.Group, error)
	GroupUpdate(g *models.Group) (*models.Group, error)
	GroupDocInsert(g *models.Group) (*models.Group, error)
}
