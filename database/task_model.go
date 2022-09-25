package database

import (
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// taskModel structures a group BSON document to save in a users collection
type taskModel struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Completed    bool               `bson:"completed,omitempty"`
	Due          time.Time          `bson:"due,omitempty"`
	Description  string             `bson:"description,omitempty"`
	UserId       primitive.ObjectID `bson:"user_id,omitempty"`
	GroupId      primitive.ObjectID `bson:"group_id,omitempty"`
	LastModified time.Time          `bson:"last_modified,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	DeletedAt    time.Time          `bson:"deleted_at,omitempty"`
}

// newTaskModel initializes a new pointer to a userModel struct from a pointer to a JSON User struct
func newTaskModel(u *models.Task) (um *taskModel, err error) {
	um = &taskModel{
		Name:         u.Name,
		Completed:    u.Completed,
		Due:          u.Due,
		Description:  u.Description,
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
	if u.UserId != "" && u.UserId != "000000000000000000000000" {
		um.UserId, err = primitive.ObjectIDFromHex(u.UserId)
	}
	return
}

// update the userModel using an overwrite bson.D doc
func (u *taskModel) update(doc interface{}) (err error) {
	data, err := bsonMarshall(doc)
	if err != nil {
		return
	}
	um := taskModel{}
	err = bson.Unmarshal(data, &um)
	if len(um.Name) > 0 {
		u.Name = um.Name
	}
	if um.Completed {
		u.Completed = um.Completed
	}
	if !um.Due.IsZero() {
		u.Due = um.Due
	}
	if len(um.Description) > 0 {
		u.Description = um.Description
	}
	if len(um.UserId.Hex()) > 0 && um.UserId.Hex() != "000000000000000000000000" {
		u.UserId = um.UserId
	}
	if len(um.GroupId.Hex()) > 0 && um.GroupId.Hex() != "000000000000000000000000" {
		u.GroupId = um.GroupId
	}
	if !um.LastModified.IsZero() {
		u.LastModified = um.LastModified
	}
	return
}

// bsonLoad loads a bson doc into the userModel
func (u *taskModel) bsonLoad(doc bson.D) (err error) {
	bData, err := bsonMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, u)
	return err
}

// match compares an input bson doc and returns whether there's a match with the userModel
// TODO: Find a better way to write these model match methods
func (u *taskModel) match(doc interface{}) bool {
	data, err := bsonMarshall(doc)
	if err != nil {
		return false
	}
	um := taskModel{}
	err = bson.Unmarshal(data, &um)
	if um.Id.Hex() != "" && um.Id.Hex() != "000000000000000000000000" {
		if u.Id == um.Id {
			return true
		}
		return false
	}
	if um.UserId.Hex() != "" && um.UserId.Hex() != "000000000000000000000000" {
		if u.UserId == um.UserId {
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
func (u *taskModel) getID() (id interface{}) {
	return u.Id
}

// addTimeStamps updates an userModel struct with a timestamp
func (u *taskModel) addTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	u.LastModified = currentTime
	if newRecord {
		u.CreatedAt = currentTime
	}
}

// addObjectID checks if a userModel has a value assigned for Id if no value a new one is generated and assigned
func (u *taskModel) addObjectID() {
	if u.Id.Hex() == "" || u.Id.Hex() == "000000000000000000000000" {
		u.Id = primitive.NewObjectID()
	}
}

// postProcess updates an userModel struct postProcess to do things such as removing the password field's value
func (u *taskModel) postProcess() (err error) {
	//u.Password = ""
	if u.UserId.Hex() == "" {
		err = errors.New("user record does not have an email")
	}
	// TODO - When implementing soft delete, DeletedAt can be checked here to ensure deleted users are filtered out
	return
}

// toDoc converts the bson userModel into a bson.D
func (u *taskModel) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(u)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// bsonFilter generates a bson filter for MongoDB queries from the userModel data
func (u *taskModel) bsonFilter() (doc bson.D, err error) {
	if u.Id.Hex() != "" && u.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", u.Id}}
	} else if u.GroupId.Hex() != "" && u.GroupId.Hex() != "000000000000000000000000" {
		doc = bson.D{{"group_id", u.GroupId}}
	} else if u.UserId.Hex() != "" && u.UserId.Hex() != "000000000000000000000000" {
		doc = bson.D{{"user_id", u.UserId}}
	}
	return
}

// bsonUpdate generates a bson update for MongoDB queries from the userModel data
func (u *taskModel) bsonUpdate() (doc bson.D, err error) {
	inner, err := u.toDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}

// toRoot creates and return a new pointer to a User JSON struct from a pointer to a BSON userModel
func (u *taskModel) toRoot() *models.Task {
	return &models.Task{
		Id:           u.Id.Hex(),
		Name:         u.Name,
		Completed:    u.Completed,
		Due:          u.Due,
		Description:  u.Description,
		UserId:       u.UserId.Hex(),
		GroupId:      u.GroupId.Hex(),
		LastModified: u.LastModified,
		CreatedAt:    u.CreatedAt,
		DeletedAt:    u.DeletedAt,
	}
}
