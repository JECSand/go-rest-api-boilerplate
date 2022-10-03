package database

import (
	"context"
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"time"
)

// TaskService is used by the app to manage all Task related controllers and functionality
type TaskService struct {
	collection   DBCollection
	db           DBClient
	taskHandler  *DBHandler[*taskModel]
	userHandler  *DBHandler[*userModel]
	groupHandler *DBHandler[*groupModel]
}

// NewTaskService is an exported function used to initialize a new TaskService struct
func NewTaskService(db DBClient, tHandler *DBHandler[*taskModel], uHandler *DBHandler[*userModel], gHandler *DBHandler[*groupModel]) *TaskService {
	collection := db.GetCollection("tasks")
	return &TaskService{collection, db, tHandler, uHandler, gHandler}
}

// checkLinkedRecords ensures the userId and groupId in the models.Task is correct
func (p *TaskService) checkLinkedRecords(g *groupModel, u *userModel) error {
	gOutCh := make(chan *groupModel)
	gErrCh := make(chan error)
	uOutCh := make(chan *userModel)
	uErrCh := make(chan error)
	go func() {
		reG, err := p.groupHandler.FindOne(g)
		gOutCh <- reG
		gErrCh <- err
	}()
	go func() {
		reU, err := p.userHandler.FindOne(u)
		uOutCh <- reU
		uErrCh <- err
	}()
	for i := 0; i < 4; i++ {
		select {
		case gOut := <-gOutCh:
			g = gOut
		case gErr := <-gErrCh:
			if gErr != nil {
				return errors.New("invalid group id")
			}
		case uOut := <-uOutCh:
			u = uOut
		case uErr := <-uErrCh:
			if uErr != nil {
				return errors.New("invalid user id")
			}
		}
	}
	if g.Id != u.GroupId {
		return errors.New("task user is not in task group")
	}
	return nil
}

// TaskCreate is used to create a new user Task
func (p *TaskService) TaskCreate(g *models.Task) (*models.Task, error) {
	err := g.Validate("create")
	if err != nil {
		return nil, err
	}
	gm, err := newTaskModel(g)
	if err != nil {
		return nil, err
	}
	err = p.checkLinkedRecords(&groupModel{Id: gm.GroupId}, &userModel{Id: gm.UserId})
	if err != nil {
		return nil, err
	}
	gm, err = p.taskHandler.InsertOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// TasksFind is used to find all Task docs in a MongoDB Collection
func (p *TaskService) TasksFind(g *models.Task) ([]*models.Task, error) {
	var tasks []*models.Task
	tm, err := newTaskModel(g)
	if err != nil {
		return tasks, err
	}
	gms, err := p.taskHandler.FindMany(tm)
	if err != nil {
		return tasks, err
	}
	for _, gm := range gms {
		tasks = append(tasks, gm.toRoot())
	}
	return tasks, nil
}

// TaskFind is used to find a specific Task doc
func (p *TaskService) TaskFind(g *models.Task) (*models.Task, error) {
	gm, err := newTaskModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.taskHandler.FindOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// TaskDelete is used to delete a Task doc
func (p *TaskService) TaskDelete(g *models.Task) (*models.Task, error) {
	gm, err := newTaskModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.taskHandler.DeleteOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// TaskDeleteMany is used to delete many Tasks
func (p *TaskService) TaskDeleteMany(g *models.Task) (*models.Task, error) {
	gm, err := newTaskModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.taskHandler.DeleteMany(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// TaskUpdate is used to update an existing Task
func (p *TaskService) TaskUpdate(g *models.Task) (*models.Task, error) {
	var filter models.Task
	err := g.Validate("update")
	if err != nil {
		return nil, err
	}
	filter.Id = g.Id
	f, err := newTaskModel(&filter)
	if err != nil {
		return nil, err
	}
	cur, TaskErr := p.taskHandler.FindOne(f)
	if TaskErr != nil {
		return nil, errors.New("task not found")
	}
	g.BuildUpdate(cur.toRoot())
	gm, err := newTaskModel(g)
	if err != nil {
		return nil, err
	}
	err = p.checkLinkedRecords(&groupModel{Id: gm.GroupId}, &userModel{Id: gm.UserId})
	if err != nil {
		return nil, err
	}
	gm, err = p.taskHandler.UpdateOne(f, gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// TaskDocInsert is used to insert a Task doc directly into mongodb for testing purposes
func (p *TaskService) TaskDocInsert(g *models.Task) (*models.Task, error) {
	insertTask, err := newTaskModel(g)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = p.collection.InsertOne(ctx, insertTask)
	if err != nil {
		return nil, err
	}
	return insertTask.toRoot(), nil
}
