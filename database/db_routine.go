package database

import (
	"sync"
)

type routineType int64

const (
	FindOne routineType = iota
	UpdateOne
	InsertOne
	DeleteOne
)

// dbRoutine struct to store executing thread ...
type dbRoutine[T dbModel] struct {
	out     T
	err     error
	handler *DBHandler[T]
	rType   routineType
	filter  T
	data    T
}

// execute a DB Routine by inputting a RoutineType, filter, and data
func (p *dbRoutine[T]) execute(rt routineType, tCh chan T, eCh chan error, f T, d T) {
	p.rType = rt
	p.filter = f
	p.data = d
	var resp T
	var err error
	switch p.rType {
	case FindOne:
		resp, err = p.handler.FindOne(p.filter)
	case UpdateOne:
		resp, err = p.handler.UpdateOne(p.filter, p.data)
	case InsertOne:
		resp, err = p.handler.InsertOne(p.data)
	case DeleteOne:
		resp, err = p.handler.DeleteOne(p.filter)
	}
	eCh <- err
	tCh <- resp
}

// resolve an executing dbRoutine
func (p *dbRoutine[T]) resolve(tCh chan T, eCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 2; i++ {
		select {
		case gOut := <-tCh:
			p.out = gOut
		case gErr := <-eCh:
			p.err = gErr
		}
	}
}
