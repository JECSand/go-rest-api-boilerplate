package database

import (
	"errors"
	"fmt"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type blacklistModel struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	AuthToken    string             `bson:"auth_token,omitempty"`
	LastModified time.Time          `bson:"last_modified,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
}

// newBlacklistModel initializes a new pointer to a blacklistModel struct from a pointer to a JSON Blacklist struct
func newBlacklistModel(bl *models.Blacklist) (bm *blacklistModel, err error) {
	bm = &blacklistModel{
		AuthToken:    bl.AuthToken,
		LastModified: bl.LastModified,
		CreatedAt:    bl.CreatedAt,
	}
	if bl.Id != "" && bl.Id != "000000000000000000000000" {
		bm.Id, err = primitive.ObjectIDFromHex(bl.Id)
	}
	return
}

// update the blacklistModel using an overwrite bson doc
func (b *blacklistModel) update(doc interface{}) (err error) {
	data, err := bsonMarshall(doc)
	if err != nil {
		return
	}
	bm := blacklistModel{}
	err = bson.Unmarshal(data, &bm)
	if len(bm.AuthToken) > 0 {
		b.AuthToken = bm.AuthToken
	}
	if !bm.LastModified.IsZero() {
		b.LastModified = bm.LastModified
	}
	return
}

// bsonLoad loads a bson doc into the blacklistModel
func (b *blacklistModel) bsonLoad(doc bson.D) (err error) {
	bData, err := bsonMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, b)
	return err
}

// match compares an input bson doc and returns whether there's a match with the blacklistModel
func (b *blacklistModel) match(doc interface{}) bool {
	data, err := bsonMarshall(doc)
	if err != nil {
		return false
	}
	bm := blacklistModel{}
	err = bson.Unmarshal(data, &bm)
	fmt.Println("\nCHECK BLACKLIST MATCH: ", bm, b)
	if b.Id == bm.Id {
		return true
	}
	if b.AuthToken == bm.AuthToken {
		return true
	}
	return false
}

// getID returns the unique identifier of the blacklistModel
func (b *blacklistModel) getID() (id interface{}) {
	return b.Id
}

// addTimeStamps updates a blacklistModel struct with a timestamp
func (b *blacklistModel) addTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	b.LastModified = currentTime
	if newRecord {
		b.CreatedAt = currentTime
	}
}

// addObjectID checks if a blacklistModel has a value assigned for Id, if no value a new one is generated and assigned
func (b *blacklistModel) addObjectID() {
	if b.Id.Hex() == "" || b.Id.Hex() == "000000000000000000000000" {
		b.Id = primitive.NewObjectID()
	}
}

// postProcess updates an blacklistModel struct postProcess to do things such as removing the password field's value
func (b *blacklistModel) postProcess() (err error) {
	if b.AuthToken == "" {
		err = errors.New("blacklist record does not have an AuthToken")
	}
	return
}

// toDoc converts the bson blacklistModel into a bson.D
func (b *blacklistModel) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(b)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// bsonFilter generates a bson filter for MongoDB queries from the blacklistModel data
func (b *blacklistModel) bsonFilter() (doc bson.D, err error) {
	if b.AuthToken != "" {
		doc = bson.D{{"auth_token", b.AuthToken}}
	} else if b.Id.Hex() != "" && b.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", b.Id}}
	}
	return
}

// bsonUpdate generates a bson update for MongoDB queries from the blacklistModel data
func (b *blacklistModel) bsonUpdate() (doc bson.D, err error) {
	inner, err := b.toDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}

// toRoot creates and return a new pointer to a Blacklist JSON struct from a pointer to a BSON blacklistModel
func (b *blacklistModel) toRoot() *models.Blacklist {
	return &models.Blacklist{
		Id:           b.Id.Hex(),
		AuthToken:    b.AuthToken,
		LastModified: b.LastModified,
		CreatedAt:    b.CreatedAt,
	}
}
