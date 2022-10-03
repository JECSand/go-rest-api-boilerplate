package services

import (
	"bytes"
	"github.com/JECSand/go-rest-api-boilerplate/models"
)

// FileService is an interface used to manage the relevant file doc controllers
type FileService interface {
	FileCreate(g *models.File, content []byte) (*models.File, error)
	FileFind(g *models.File) (*models.File, error)
	FilesFind(g *models.File) ([]*models.File, error)
	FileDelete(g *models.File) (*models.File, error)
	FileDeleteMany(g []*models.File) error
	FileUpdate(g *models.File, content []byte) (*models.File, error)
	RetrieveFile(g *models.File) (*bytes.Buffer, error)
}
