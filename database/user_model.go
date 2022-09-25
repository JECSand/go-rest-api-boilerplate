package database

import (
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// userModel structures a group BSON document to save in a users collection
type userModel struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username,omitempty"`
	Password     string             `bson:"password,omitempty"`
	FirstName    string             `bson:"firstname,omitempty"`
	LastName     string             `bson:"lastname,omitempty"`
	Email        string             `bson:"email,omitempty"`
	Role         string             `bson:"role,omitempty"`
	RootAdmin    bool               `bson:"root_admin,omitempty"`
	GroupId      primitive.ObjectID `bson:"group_id,omitempty"`
	LastModified time.Time          `bson:"last_modified,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	DeletedAt    time.Time          `bson:"deleted_at,omitempty"`
}

// newUserModel initializes a new pointer to a userModel struct from a pointer to a JSON User struct
func newUserModel(u *models.User) (um *userModel, err error) {
	um = &userModel{
		Username:     u.Username,
		Password:     u.Password,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		Role:         u.Role,
		RootAdmin:    u.RootAdmin,
		LastModified: u.LastModified,
		CreatedAt:    u.CreatedAt,
		DeletedAt:    u.DeletedAt,
	}
	if u.Id != "" && u.Id != "000000000000000000000000" {
		um.Id, err = primitive.ObjectIDFromHex(u.Id)
	}
	if u.GroupId != "" && u.GroupId != "000000000000000000000000" {
		um.GroupId, err = primitive.ObjectIDFromHex(u.GroupId)
	}
	return
}

// update the userModel using an overwrite bson.D doc
func (u *userModel) update(doc interface{}) (err error) {
	data, err := bsonMarshall(doc)
	if err != nil {
		return
	}
	um := userModel{}
	err = bson.Unmarshal(data, &um)
	if len(um.Username) > 0 {
		u.Username = um.Username
	}
	if len(um.FirstName) > 0 {
		u.FirstName = um.FirstName
	}
	if len(um.LastName) > 0 {
		u.LastName = um.LastName
	}
	if len(um.Email) > 0 {
		u.Email = um.Email
	}
	if len(um.Password) > 0 {
		u.Password = um.Password
	}
	if len(um.GroupId.Hex()) > 0 && um.GroupId.Hex() != "000000000000000000000000" {
		u.GroupId = um.GroupId
	}
	if len(um.Role) > 0 {
		u.Role = um.Role
	}
	if !um.LastModified.IsZero() {
		u.LastModified = um.LastModified
	}
	return
}

// bsonLoad loads a bson doc into the userModel
func (u *userModel) bsonLoad(doc bson.D) (err error) {
	bData, err := bsonMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, u)
	return err
}

// match compares an input bson doc and returns whether there's a match with the userModel
// TODO: Find a better way to write these model match methods
func (u *userModel) match(doc interface{}) bool {
	data, err := bsonMarshall(doc)
	if err != nil {
		return false
	}
	um := userModel{}
	err = bson.Unmarshal(data, &um)
	if um.Id.Hex() != "" && um.Id.Hex() != "000000000000000000000000" {
		if u.Id == um.Id {
			return true
		}
		return false
	}
	if um.Email != "" {
		if u.Email == um.Email {
			return true
		}
		return false
	}
	if um.GroupId.Hex() != "" && um.GroupId.Hex() != "000000000000000000000000" {
		if u.GroupId == um.GroupId {
			return true
		}
		return false
	}
	return false
}

// getID returns the unique identifier of the userModel
func (u *userModel) getID() (id interface{}) {
	return u.Id
}

// addTimeStamps updates an userModel struct with a timestamp
func (u *userModel) addTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	u.LastModified = currentTime
	if newRecord {
		u.CreatedAt = currentTime
	}
}

// addObjectID checks if a userModel has a value assigned for Id if no value a new one is generated and assigned
func (u *userModel) addObjectID() {
	if u.Id.Hex() == "" || u.Id.Hex() == "000000000000000000000000" {
		u.Id = primitive.NewObjectID()
	}
}

// postProcess updates an userModel struct postProcess to do things such as removing the password field's value
func (u *userModel) postProcess() (err error) {
	//u.Password = ""
	if u.Email == "" {
		err = errors.New("user record does not have an email")
	}
	// TODO - When implementing soft delete, DeletedAt can be checked here to ensure deleted users are filtered out
	return
}

// toDoc converts the bson userModel into a bson.D
func (u *userModel) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(u)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// bsonFilter generates a bson filter for MongoDB queries from the userModel data
func (u *userModel) bsonFilter() (doc bson.D, err error) {
	if u.Id.Hex() != "" && u.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", u.Id}}
	} else if u.GroupId.Hex() != "" && u.GroupId.Hex() != "000000000000000000000000" {
		doc = bson.D{{"group_id", u.GroupId}}
	} else if u.Email != "" {
		doc = bson.D{{"email", u.Email}}
	}
	return
}

// bsonUpdate generates a bson update for MongoDB queries from the userModel data
func (u *userModel) bsonUpdate() (doc bson.D, err error) {
	inner, err := u.toDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}

// toRoot creates and return a new pointer to a User JSON struct from a pointer to a BSON userModel
func (u *userModel) toRoot() *models.User {
	return &models.User{
		Id:           u.Id.Hex(),
		Username:     u.Username,
		Password:     u.Password,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		Role:         u.Role,
		RootAdmin:    u.RootAdmin,
		GroupId:      u.GroupId.Hex(),
		LastModified: u.LastModified,
		CreatedAt:    u.CreatedAt,
		DeletedAt:    u.DeletedAt,
	}
}
