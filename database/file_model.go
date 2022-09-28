package database

import (
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// fileModel structures a file BSON document associated with a GridFS Object to save in a files collection
type fileModel struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId      primitive.ObjectID `bson:"owner_id,omitempty"`
	OwnerType    string             `bson:"owner_type,omitempty"`
	GridFSId     primitive.ObjectID `bson:"gridfs_id,omitempty"`
	BucketName   string             `bson:"bucket_name,omitempty"`
	BucketType   string             `bson:"bucket_type,omitempty"`
	Name         string             `bson:"name,omitempty"`
	FileType     string             `bson:"file_type,omitempty"`
	Size         int                `bson:"size,omitempty"`
	LastModified time.Time          `bson:"last_modified,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	DeletedAt    time.Time          `bson:"deleted_at,omitempty"`
}

// newFileModel initializes a new pointer to a fileModel struct from a pointer to a JSON User struct
func newFileModel(u *models.File) (um *fileModel, err error) {
	um = &fileModel{
		OwnerType:    u.OwnerType,
		BucketName:   u.BucketName,
		BucketType:   u.BucketType,
		Name:         u.Name,
		FileType:     u.FileType,
		Size:         u.Size,
		LastModified: u.LastModified,
		CreatedAt:    u.CreatedAt,
		DeletedAt:    u.DeletedAt,
	}
	if u.Id != "" && u.Id != "000000000000000000000000" {
		um.Id, err = primitive.ObjectIDFromHex(u.Id)
	}
	if u.OwnerId != "" && u.OwnerId != "000000000000000000000000" {
		um.OwnerId, err = primitive.ObjectIDFromHex(u.OwnerId)
	}
	if u.GridFSId != "" && u.GridFSId != "000000000000000000000000" {
		um.GridFSId, err = primitive.ObjectIDFromHex(u.GridFSId)
	}
	return
}

// update the userModel using an overwrite bson.D doc
func (u *fileModel) update(doc interface{}) (err error) {
	data, err := bsonMarshall(doc)
	if err != nil {
		return
	}
	um := fileModel{}
	err = bson.Unmarshal(data, &um)
	if len(um.Name) > 0 {
		u.Name = um.Name
	}
	if len(um.BucketName) > 0 {
		u.BucketName = um.BucketName
	}
	if len(um.BucketType) > 0 {
		u.BucketType = um.BucketType
	}
	if um.Size > 0 {
		u.Size = um.Size
	}
	if len(um.OwnerId.Hex()) > 0 && um.OwnerId.Hex() != "000000000000000000000000" {
		u.OwnerId = um.OwnerId
	}
	if len(um.GridFSId.Hex()) > 0 && um.GridFSId.Hex() != "000000000000000000000000" {
		u.GridFSId = um.GridFSId
	}
	if !um.LastModified.IsZero() {
		u.LastModified = um.LastModified
	}
	return
}

// bsonLoad loads a bson doc into the userModel
func (u *fileModel) bsonLoad(doc bson.D) (err error) {
	bData, err := bsonMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, u)
	return err
}

// match compares an input bson doc and returns whether there's a match with the userModel
// TODO: Find a better way to write these model match methods
func (u *fileModel) match(doc interface{}) bool {
	data, err := bsonMarshall(doc)
	if err != nil {
		return false
	}
	um := fileModel{}
	err = bson.Unmarshal(data, &um)
	if um.Id.Hex() != "" && um.Id.Hex() != "000000000000000000000000" {
		if u.Id == um.Id {
			return true
		}
		return false
	}
	if um.OwnerId.Hex() != "" && um.OwnerId.Hex() != "000000000000000000000000" {
		if u.OwnerId == um.OwnerId {
			return true
		}
		return false
	}
	if um.GridFSId.Hex() != "" && um.GridFSId.Hex() != "000000000000000000000000" {
		if u.GridFSId == um.GridFSId {
			return true
		}
		return false
	}
	return false
}

// getID returns the unique identifier of the userModel
func (u *fileModel) getID() (id interface{}) {
	return u.Id
}

// addTimeStamps updates an userModel struct with a timestamp
func (u *fileModel) addTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	u.LastModified = currentTime
	if newRecord {
		u.CreatedAt = currentTime
	}
}

// addObjectID checks if a userModel has a value assigned for Id if no value a new one is generated and assigned
func (u *fileModel) addObjectID() {
	if u.Id.Hex() == "" || u.Id.Hex() == "000000000000000000000000" {
		u.Id = primitive.NewObjectID()
	}
}

// postProcess updates an userModel struct postProcess to do things such as removing the password field's value
func (u *fileModel) postProcess() (err error) {
	if u.GridFSId.Hex() == "" {
		err = errors.New("user record does not have an email")
	}
	// TODO - When implementing soft delete, DeletedAt can be checked here to ensure deleted users are filtered out
	return
}

// toDoc converts the bson userModel into a bson.D
func (u *fileModel) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(u)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// bsonFilter generates a bson filter for MongoDB queries from the userModel data
func (u *fileModel) bsonFilter() (doc bson.D, err error) {
	if u.Id.Hex() != "" && u.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", u.Id}}
	} else if u.OwnerId.Hex() != "" && u.OwnerId.Hex() != "000000000000000000000000" {
		doc = bson.D{{"owner_id", u.OwnerId}}
	} else if u.GridFSId.Hex() != "" && u.GridFSId.Hex() != "000000000000000000000000" {
		doc = bson.D{{"gridfs_id", u.GridFSId}}
	}
	return
}

// bsonUpdate generates a bson update for MongoDB queries from the userModel data
func (u *fileModel) bsonUpdate() (doc bson.D, err error) {
	inner, err := u.toDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}

// toRoot creates and return a new pointer to a File JSON struct from a pointer to a BSON fileModel
func (u *fileModel) toRoot() *models.File {
	return &models.File{
		Id:           u.Id.Hex(),
		OwnerId:      u.OwnerId.Hex(),
		OwnerType:    u.OwnerType,
		GridFSId:     u.GridFSId.Hex(),
		BucketName:   u.BucketName,
		BucketType:   u.BucketType,
		Name:         u.Name,
		FileType:     u.FileType,
		Size:         u.Size,
		LastModified: u.LastModified,
		CreatedAt:    u.CreatedAt,
		DeletedAt:    u.DeletedAt,
	}
}
