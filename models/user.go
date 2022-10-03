package models

import (
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/utilities"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// User is a root struct that is used to store the json encoded data for/from a mongodb user doc.
type User struct {
	Id           string    `json:"id,omitempty"`
	Username     string    `json:"username,omitempty"`
	Password     string    `json:"password,omitempty"`
	FirstName    string    `json:"firstname,omitempty"`
	LastName     string    `json:"lastname,omitempty"`
	Email        string    `json:"email,omitempty"`
	Role         string    `json:"role,omitempty"`
	RootAdmin    bool      `json:"root_admin,omitempty"`
	GroupId      string    `json:"group_id,omitempty"`
	ImageId      string    `json:"image_id,omitempty"`
	LastModified time.Time `json:"last_modified,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	DeletedAt    time.Time `json:"deleted_at,omitempty"`
}

// LoadScope scopes the User struct
func (g *User) LoadScope(scopeUser *User, valCase string) {
	switch valCase {
	case "create":
		if !scopeUser.RootAdmin {
			g.RootAdmin = false
			g.GroupId = scopeUser.GroupId
		}
	case "update":
		g.Id = scopeUser.Id
		if !scopeUser.RootAdmin {
			g.RootAdmin = false
			g.GroupId = scopeUser.GroupId
			if scopeUser.Role != "admin" {
				g.Role = "member"
			}
		}
	case "find":
		if !scopeUser.RootAdmin {
			g.GroupId = scopeUser.GroupId
			if scopeUser.Role != "admin" {
				g.Id = scopeUser.Id
			}
		}
	}
	return
}

// CheckID determines whether a specified ID is set or not
func (g *User) CheckID(chkId string) bool {
	switch chkId {
	case "id":
		if !utilities.CheckObjectID(g.Id) {
			return false
		}
	case "group_id":
		if !utilities.CheckObjectID(g.GroupId) {
			return false
		}
	case "image_id":
		if !utilities.CheckObjectID(g.ImageId) {
			return false
		}
	}
	return true
}

// Authenticate compares an input password with the hashed password stored in the User model
func (g *User) Authenticate(checkPassword string) error {
	if len(g.Password) != 0 {
		password := []byte(g.Password)
		cPassword := []byte(checkPassword)
		return bcrypt.CompareHashAndPassword(password, cPassword)
	}
	return errors.New("no password set to hash in user model")
}

// HashPassword hashes a user password and associates it with the user struct
func (g *User) HashPassword() error {
	if len(g.Password) != 0 {
		password := []byte(g.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		g.Password = string(hashedPassword)
		return nil
	}
	return errors.New("no password set to hash in user model")
}

// Validate a User for different scenarios such as loading TokenData, creating new User, or updating a User
func (g *User) Validate(valCase string) (err error) {
	var missingFields []string
	switch valCase {
	case "auth":
		if !g.CheckID("id") {
			missingFields = append(missingFields, "id")
		}
		if !g.CheckID("group_id") {
			missingFields = append(missingFields, "group_id")
		}
		if g.Role == "" {
			missingFields = append(missingFields, "role")
		}
	case "create":
		if g.Username == "" {
			missingFields = append(missingFields, "id")
		}
		if g.Email == "" {
			missingFields = append(missingFields, "email")
		}
		if g.Password == "" {
			missingFields = append(missingFields, "password")
		}
		if !g.CheckID("group_id") {
			missingFields = append(missingFields, "group_id")
		}
	case "update":
		if !g.CheckID("id") && g.Email == "" {
			missingFields = append(missingFields, "id")
		}
	default:
		return errors.New("unrecognized validation case")
	}
	if len(missingFields) > 0 {
		return errors.New("missing the following user fields: " + strings.Join(missingFields, ", "))
	}
	return
}

// BuildFilter is a function that setups the base user struct during a user modification request
func (g *User) BuildFilter() (*User, error) {
	var filter User
	if g.CheckID("id") {
		filter.Id = g.Id
	} else if g.Email != "" {
		filter.Email = g.Email
	} else {
		return nil, errors.New("user is missing a valid query filter")
	}
	return &filter, nil
}

// BuildUpdate is a function that setups the base user struct during a user modification request
func (g *User) BuildUpdate(curUser *User) {
	if len(g.Username) == 0 {
		g.Username = curUser.Username
	}
	if len(g.FirstName) == 0 {
		g.FirstName = curUser.FirstName
	}
	if len(g.LastName) == 0 {
		g.LastName = curUser.LastName
	}
	if len(g.Email) == 0 {
		g.Email = curUser.Email
	}
	if len(g.GroupId) == 0 {
		g.GroupId = curUser.GroupId
	}
	if len(g.ImageId) == 0 {
		g.ImageId = curUser.ImageId
	}
	if len(g.Role) == 0 {
		g.Role = curUser.Role
	}
}

// UsersToFiles converts an input slice of user to a slice of file
func UsersToFiles(users []*User) []*File {
	var files []*File
	for _, u := range users {
		if u.CheckID("image_id") {
			files = append(files, &File{OwnerId: u.Id, OwnerType: "user"})
		}
	}
	return files
}
