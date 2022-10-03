package models

import (
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/utilities"
	"strings"
	"time"
)

// Task is a root struct that is used to store the json encoded data for/from a mongodb group doc.
type Task struct {
	Id           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Completed    bool      `json:"completed,omitempty"`
	Due          time.Time `json:"due,omitempty"`
	Description  string    `json:"description,omitempty"`
	UserId       string    `json:"user_id,omitempty"`
	GroupId      string    `json:"group_id,omitempty"`
	LastModified time.Time `json:"last_modified,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	DeletedAt    time.Time `json:"deleted_at,omitempty"`
}

// LoadScope scopes the Task struct
func (g *Task) LoadScope(scopeUser *User) {
	if !scopeUser.RootAdmin {
		g.GroupId = scopeUser.GroupId
		if scopeUser.Role != "admin" {
			g.UserId = scopeUser.Id
		}
	}
	if !g.CheckID("user_id") {
		g.UserId = scopeUser.Id
	}
	if !g.CheckID("group_id") {
		g.GroupId = scopeUser.GroupId
	}
	return
}

// CheckID determines whether a specified ID is set or not
func (g *Task) CheckID(chkId string) bool {
	switch chkId {
	case "id":
		if !utilities.CheckObjectID(g.Id) {
			return false
		}
	case "group_id":
		if !utilities.CheckObjectID(g.GroupId) {
			return false
		}
	case "user_id":
		if !utilities.CheckObjectID(g.UserId) {
			return false
		}
	}
	return true
}

// Validate a Group for different scenarios such as loading TokenData, creating new Group, or updating a Group
func (g *Task) Validate(valCase string) (err error) {
	var missingFields []string
	switch valCase {
	case "create":
		if g.Name == "" {
			missingFields = append(missingFields, "name")
		}
		if !g.CheckID("user_id") {
			missingFields = append(missingFields, "user_id")
		}
		if !g.CheckID("group_id") {
			missingFields = append(missingFields, "group_id")
		}
		if g.Due.IsZero() {
			missingFields = append(missingFields, "due")
		}
	case "update":
		if !g.CheckID("id") {
			missingFields = append(missingFields, "id")
		}
	default:
		return errors.New("unrecognized validation case")
	}
	if len(missingFields) > 0 {
		return errors.New("missing the following group fields: " + strings.Join(missingFields, ", "))
	}
	return
}

// BuildUpdate is a function that setups the base task struct during a user modification request
func (g *Task) BuildUpdate(cur *Task) {
	if len(g.Name) == 0 {
		g.Name = cur.Name
	}
	if !g.Completed {
		g.Completed = cur.Completed
	}
	if g.Due.IsZero() {
		g.Due = cur.Due
	}
	if len(g.Description) == 0 {
		g.Description = cur.Description
	}
	if len(g.UserId) == 0 {
		g.UserId = cur.UserId
	}
	if len(g.GroupId) == 0 {
		g.GroupId = cur.GroupId
	}
}
