package database

import (
	"context"
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"time"
)

// GroupService is used by the app to manage all group related controllers and functionality
type GroupService struct {
	collection DBCollection
	db         DBClient
	handler    *DBHandler[*groupModel]
}

// NewGroupService is an exported function used to initialize a new GroupService struct
func NewGroupService(db DBClient, handler *DBHandler[*groupModel]) *GroupService {
	collection := db.GetCollection("groups")
	return &GroupService{collection, db, handler}
}

// GroupCreate is used to create a new user group
func (p *GroupService) GroupCreate(g *models.Group) (*models.Group, error) {
	err := g.Validate("create")
	if err != nil {
		return nil, err
	}
	gm, err := newGroupModel(g)
	if err != nil {
		return nil, err
	}
	_, err = p.handler.FindOne(&groupModel{Name: gm.Name})
	if err == nil {
		return nil, errors.New("group name exists")
	}
	gm, err = p.handler.InsertOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// GroupsFind is used to find all group docs in a MongoDB Collection
func (p *GroupService) GroupsFind(g *models.Group) ([]*models.Group, error) {
	var groups []*models.Group
	m, err := newGroupModel(g)
	if err != nil {
		return groups, err
	}
	gms, err := p.handler.FindMany(m)
	if err != nil {
		return groups, err
	}
	for _, gm := range gms {
		groups = append(groups, gm.toRoot())
	}
	return groups, nil
}

// GroupFind is used to find a specific group doc
func (p *GroupService) GroupFind(g *models.Group) (*models.Group, error) {
	gm, err := newGroupModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.handler.FindOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// GroupDelete is used to delete a group doc
func (p *GroupService) GroupDelete(g *models.Group) (*models.Group, error) {
	gm, err := newGroupModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.handler.DeleteOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// GroupDeleteMany is used to delete many Groups
func (p *GroupService) GroupDeleteMany(g *models.Group) (*models.Group, error) {
	gm, err := newGroupModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.handler.DeleteMany(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// GroupUpdate is used to update an existing group
func (p *GroupService) GroupUpdate(g *models.Group) (*models.Group, error) {
	var filter models.Group
	err := g.Validate("create")
	if err != nil {
		return nil, errors.New("missing valid query filter")
	}
	filter.Id = g.Id
	if g.Name != "" {
		reDoc, err := p.handler.FindOne(&groupModel{Name: g.Name})
		if err == nil && reDoc.toRoot().Id != filter.Id {
			return nil, errors.New("group name exists")
		}
	}
	f, err := newGroupModel(&filter)
	if err != nil {
		return nil, err
	}
	gm, err := newGroupModel(g)
	if err != nil {
		return nil, err
	}
	_, groupErr := p.handler.FindOne(f)
	if groupErr != nil {
		return nil, errors.New("group not found")
	}
	gm, err = p.handler.UpdateOne(f, gm)
	return gm.toRoot(), err
}

// GroupDocInsert is used to insert a group doc directly into mongodb for testing purposes
func (p *GroupService) GroupDocInsert(g *models.Group) (*models.Group, error) {
	insertGroup, err := newGroupModel(g)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = p.collection.InsertOne(ctx, insertGroup)
	if err != nil {
		return nil, err
	}
	return insertGroup.toRoot(), nil
}
