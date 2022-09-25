package models

import (
	"errors"
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

// checkID determines whether a specified ID is set or not
func (g *Task) checkID(chkId string) bool {
	switch chkId {
	case "id":
		if g.Id == "" || g.Id == "000000000000000000000000" {
			return false
		}
	case "group_id":
		if g.GroupId == "" || g.GroupId == "000000000000000000000000" {
			return false
		}
	case "user_id":
		if g.UserId == "" || g.UserId == "000000000000000000000000" {
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
		if !g.checkID("user_id") {
			missingFields = append(missingFields, "user_id")
		}
		if !g.checkID("group_id") {
			missingFields = append(missingFields, "group_id")
		}
		if g.Due.IsZero() {
			missingFields = append(missingFields, "due")
		}
	case "update":
		if !g.checkID("id") {
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
