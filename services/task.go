package services

import "github.com/JECSand/go-rest-api-boilerplate/models"

// TaskService is an interface used to manage the relevant group doc controllers
type TaskService interface {
	TaskCreate(g *models.Task) (*models.Task, error)
	TaskFind(g *models.Task) (*models.Task, error)
	TasksFind(g *models.Task) ([]*models.Task, error)
	TaskDelete(g *models.Task) (*models.Task, error)
	TaskUpdate(g *models.Task) (*models.Task, error)
	TaskDocInsert(g *models.Task) (*models.Task, error)
}
