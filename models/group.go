package models

import (
	"errors"
	"strings"
	"time"
)

// Group is a root struct that is used to store the json encoded data for/from a mongodb group doc.
type Group struct {
	Id           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	RootAdmin    bool      `json:"root_admin,omitempty"`
	LastModified time.Time `json:"last_modified,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	DeletedAt    time.Time `json:"deleted_at,omitempty"`
}

// checkID determines whether a specified ID is set or not
func (g *Group) checkID(chkId string) bool {
	switch chkId {
	case "id":
		if g.Id == "" || g.Id == "000000000000000000000000" {
			return false
		}
	}
	return true
}

// Validate a Group for different scenarios such as loading TokenData, creating new Group, or updating a Group
func (g *Group) Validate(valCase string) (err error) {
	var missingFields []string
	switch valCase {
	case "create":
		if g.Name == "" {
			missingFields = append(missingFields, "name")
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
