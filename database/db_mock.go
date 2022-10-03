package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"os"
	"time"
)

/*
================ testDBUtils ==================
*/

// cleanUpdateBSON inputs a bson type and attempts to marshall it into a slice of bytes
func cleanUpdateBSON(bsonData interface{}) (data interface{}, err error) {
	switch t := bsonData.(type) {
	case nil:
		return nil, errors.New("input bsonData to marshall can not be nil")
	case bson.D:
		if len(t) > 0 && t[0].Key == "$set" {
			return t[0].Value, nil
		}
		return t, nil
	}
	return bsonData, nil
}

// standardizeID ensures that a dbModels unique identified is returned as a string
func standardizeID(dbDoc dbModel) (string, error) {
	var docId string
	switch t := dbDoc.getID().(type) {
	case nil:
		return docId, errors.New("a test record being inserted is missing a unique identifier")
	case primitive.ObjectID:
		docId = t.Hex()
	case string:
		docId = t
	}
	return docId, nil
}

// bsonMarshall inputs a bson type and attempts to marshall it into a slice of bytes
func bsonMarshall(bsonData interface{}) (data []byte, err error) {
	switch t := bsonData.(type) {
	case nil:
		return nil, errors.New("input bsonData to marshall can not be nil")
	case []byte:
		return t, nil
	case bson.D:
		data, err = bson.Marshal(t)
		if err != nil {
			return
		}
	case bson.M:
		data, err = bson.Marshal(t)
		if err != nil {
			return
		}
	case bson.Raw:
		data, err = bson.Marshal(t)
		if err != nil {
			return
		}
	case bson.E:
		data, err = bson.Marshal(t)
		if err != nil {
			return
		}
	}
	return
}

// bsonUnmarshall inputs a bson type and attempts to marshall it into a slice of bytes
func bsonUnmarshall(colName string, bsonData interface{}) (data dbModel, err error) {
	switch colName {
	case "users":
		bData, err := bsonMarshall(bsonData)
		if err != nil {
			panic(err)
			return nil, err
		}
		um := userModel{}
		err = bson.Unmarshal(bData, &um)
		return &um, nil
	case "groups":
		bData, err := bsonMarshall(bsonData)
		if err != nil {
			return nil, err
		}
		gm := groupModel{}
		err = bson.Unmarshal(bData, &gm)
		return &gm, nil
	case "blacklists":
		bData, err := bsonMarshall(bsonData)
		if err != nil {
			return nil, err
		}
		bm := blacklistModel{}
		err = bson.Unmarshal(bData, &bm)
		return &bm, nil
	case "tasks":
		bData, err := bsonMarshall(bsonData)
		if err != nil {
			return nil, err
		}
		tm := taskModel{}
		err = bson.Unmarshal(bData, &tm)
		return &tm, nil
	}
	return nil, errors.New("invalid test collection type")
}

/*
================ Test Data ==================
*/

func getTestGroupModels(root bool) []*groupModel {
	var gms []*groupModel
	var gm *groupModel
	if root {
		gm, _ = newGroupModel(&models.Group{
			Id:        "000000000000000000000001",
			Name:      "test1",
			RootAdmin: true,
		})
		gms = append(gms, gm)
	}
	gm, _ = newGroupModel(&models.Group{
		Id:        "000000000000000000000002",
		Name:      "test2",
		RootAdmin: false,
	})
	gms = append(gms, gm)
	gm, _ = newGroupModel(&models.Group{
		Id:        "000000000000000000000003",
		Name:      "test3",
		RootAdmin: false,
	})
	gms = append(gms, gm)
	return gms
}

func getTestUsersModels(root bool) []*userModel {
	var gms []*userModel
	var gm *userModel
	if root {
		gm, _ = newUserModel(&models.User{
			Id:        "000000000000000000000011",
			Email:     "test1@email.com",
			Password:  "abc123",
			GroupId:   "000000000000000000000001",
			Role:      "admin",
			RootAdmin: true,
		})
		gms = append(gms, gm)
	}
	gm, _ = newUserModel(&models.User{
		Id:        "000000000000000000000012",
		Email:     "test2@email.com",
		Password:  "abc123",
		GroupId:   "000000000000000000000002",
		Role:      "member",
		RootAdmin: false,
	})
	gms = append(gms, gm)
	gm, _ = newUserModel(&models.User{
		Id:        "000000000000000000000013",
		Email:     "test3@email.com",
		Password:  "abc123",
		GroupId:   "000000000000000000000002",
		Role:      "member",
		RootAdmin: false,
	})
	gms = append(gms, gm)
	return gms
}

func getTestTasksModels() []*taskModel {
	var gms []*taskModel
	var gm *taskModel
	gm, _ = newTaskModel(&models.Task{
		Id:      "000000000000000000000022",
		Name:    "Task1",
		Due:     time.Now().UTC(),
		UserId:  "000000000000000000000013",
		GroupId: "000000000000000000000002",
	})
	gms = append(gms, gm)
	gm, _ = newTaskModel(&models.Task{
		Id:      "000000000000000000000023",
		Name:    "Task2",
		Due:     time.Now().UTC(),
		UserId:  "000000000000000000000012",
		GroupId: "000000000000000000000002",
	})
	gms = append(gms, gm)
	return gms
}

func getTestTokens() []string {
	return []string{
		"123445608654321",
		"123445678654321",
	}
}

/*
================ testTasksUtils ==================
*/

func initTestTaskService() *TaskService {
	os.Setenv("ENV", "test")
	os.Setenv("MONGO_URI", "mongodb+srv://in_mem")
	os.Setenv("DATABASE", "test")
	db, _ := initializeNewTestClient()
	gCollection := db.GetCollection("groups")
	gHandler := db.NewGroupHandler()
	gs := &GroupService{
		gCollection,
		db,
		gHandler,
	}
	tg := getTestGroupModels(true)
	for _, d := range tg {
		_, err := gs.GroupCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	uCollection := db.GetCollection("users")
	uHandler := db.NewUserHandler()
	us := &UserService{
		uCollection,
		db,
		uHandler,
		gHandler,
	}
	tu := getTestUsersModels(true)
	for _, d := range tu {
		_, err := us.UserCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	collection := db.GetCollection("tasks")
	tHandler := db.NewTaskHandler()
	return &TaskService{
		collection,
		db,
		tHandler,
		uHandler,
		gHandler,
	}
}

func setupTestTasks() *TaskService {
	os.Setenv("ENV", "test")
	os.Setenv("MONGO_URI", "mongodb+srv://in_mem")
	os.Setenv("DATABASE", "test")
	db, _ := initializeNewTestClient()
	gCollection := db.GetCollection("groups")
	gHandler := db.NewGroupHandler()
	gs := &GroupService{
		gCollection,
		db,
		gHandler,
	}
	tg := getTestGroupModels(true)
	for _, d := range tg {
		_, err := gs.GroupCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	uCollection := db.GetCollection("users")
	uHandler := db.NewUserHandler()
	us := &UserService{
		uCollection,
		db,
		uHandler,
		gHandler,
	}
	tu := getTestUsersModels(true)
	for _, d := range tu {
		_, err := us.UserCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	collection := db.GetCollection("tasks")
	tHandler := db.NewTaskHandler()
	ts := &TaskService{
		collection,
		db,
		tHandler,
		uHandler,
		gHandler,
	}
	td := getTestTasksModels()
	for _, d := range td {
		_, err := ts.TaskCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	return ts
}

/*
================ testBlacklistUtils ==================
*/

func initTestBlacklistService() *BlacklistService {
	os.Setenv("ENV", "test")
	os.Setenv("MONGO_URI", "mongodb+srv://in_mem")
	os.Setenv("DATABASE", "test")
	db, _ := initializeNewTestClient()
	collection := db.GetCollection("blacklists")
	gHandler := db.NewBlacklistHandler()
	return &BlacklistService{
		collection,
		db,
		gHandler,
	}
}

func setupTestBlacklists() *BlacklistService {
	os.Setenv("ENV", "test")
	os.Setenv("MONGO_URI", "mongodb+srv://in_mem")
	os.Setenv("DATABASE", "test")
	db, _ := initializeNewTestClient()
	collection := db.GetCollection("blacklists")
	gHandler := db.NewBlacklistHandler()
	gs := &BlacklistService{
		collection,
		db,
		gHandler,
	}
	td := getTestTokens()
	for _, d := range td {
		err := gs.BlacklistAuthToken(d)
		if err != nil {
			panic(err)
		}
	}
	return gs
}

/*
================ testGroupsUtils ==================
*/

func initTestGroupService() *GroupService {
	os.Setenv("ENV", "test")
	os.Setenv("MONGO_URI", "mongodb+srv://in_mem")
	os.Setenv("DATABASE", "test")
	db, _ := initializeNewTestClient()
	collection := db.GetCollection("groups")
	gHandler := db.NewGroupHandler()
	return &GroupService{
		collection,
		db,
		gHandler,
	}
}

func setupTestGroups() *GroupService {
	os.Setenv("ENV", "test")
	os.Setenv("MONGO_URI", "mongodb+srv://in_mem")
	os.Setenv("DATABASE", "test")
	db, _ := initializeNewTestClient()
	collection := db.GetCollection("groups")
	gHandler := db.NewGroupHandler()
	gs := &GroupService{
		collection,
		db,
		gHandler,
	}
	td := getTestGroupModels(false)
	for _, d := range td {
		_, err := gs.GroupCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	return gs
}

/*
================ testUsersUtils ==================
*/

func initTestUserService() *UserService {
	os.Setenv("ENV", "test")
	os.Setenv("MONGO_URI", "mongodb+srv://in_mem")
	os.Setenv("DATABASE", "test")
	os.Setenv("TOKEN_SECRET", "SECRET")
	db, _ := initializeNewTestClient()
	gCollection := db.GetCollection("groups")
	gHandler := db.NewGroupHandler()
	gs := &GroupService{
		gCollection,
		db,
		gHandler,
	}
	td := getTestGroupModels(true)
	for _, d := range td {
		_, err := gs.GroupCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	uCollection := db.GetCollection("users")
	uHandler := db.NewUserHandler()
	return &UserService{
		uCollection,
		db,
		uHandler,
		gHandler,
	}
}

func setupTestUsers() *UserService {
	os.Setenv("ENV", "test")
	os.Setenv("MONGO_URI", "mongodb+srv://in_mem")
	os.Setenv("DATABASE", "test")
	os.Setenv("TOKEN_SECRET", "SECRET")
	db, _ := initializeNewTestClient()
	gCollection := db.GetCollection("groups")
	gHandler := db.NewGroupHandler()
	gs := &GroupService{
		gCollection,
		db,
		gHandler,
	}
	tg := getTestGroupModels(true)
	for _, d := range tg {
		_, err := gs.GroupCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	collection := db.GetCollection("users")
	uHandler := db.NewUserHandler()
	us := &UserService{
		collection,
		db,
		uHandler,
		gHandler,
	}
	tu := getTestUsersModels(true)
	for _, d := range tu {
		_, err := us.UserCreate(d.toRoot())
		if err != nil {
			panic(err)
		}
	}
	return us
}

/*
================ testCursorData ==================
*/

// testCursorData
type testCursorData struct {
	Results []bson.D `bson:"results,omitempty"`
}

// initTestCursorData instantiates a new testCursorData
func initTestCursorData(res []dbModel) *testCursorData {
	var resBSON []bson.D
	for _, r := range res {
		d, err := r.toDoc()
		if err != nil {
			panic(err)
		}
		resBSON = append(resBSON, d)
	}
	return &testCursorData{Results: resBSON}
}

// toDoc converts the bson testCursorData into a bson.D
func (b *testCursorData) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(b)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

/*
================ testMongoCursor ==================
*/

// testMongoCursor ...
type testMongoCursor struct {
	ctx      context.Context
	Results  []byte
	docs     []bson.D
	curCurse int
}

// newTestMongoCursor initiates and returns a testMongoCursor
func newTestMongoCursor(cur *mongo.Cursor) *testMongoCursor {
	var cd testCursorData
	data, err := bsonMarshall(cur.Current)
	if err != nil {
		panic(err)
	}
	err = bson.Unmarshal(data, &cd)
	if err != nil {
		panic(err)
	}
	return &testMongoCursor{Results: cur.Current, docs: cd.Results, curCurse: 0}
}

// Next check if there's more result documents to decode
func (c *testMongoCursor) Next(ctx context.Context) bool {
	c.ctx = ctx
	if c.curCurse < len(c.docs) {
		return true
	}
	return false
}

// Decode a result document into the input val
func (c *testMongoCursor) Decode(val interface{}) error {
	if c.curCurse >= len(c.docs) {
		return errors.New("test cursor out of range")
	}
	curDoc := c.docs[c.curCurse]
	bData, err := bsonMarshall(curDoc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, val)
	if err != nil {
		return err
	}
	c.curCurse++
	return nil
}

// Close the test cursor
func (c *testMongoCursor) Close(ctx context.Context) error {
	c.ctx = ctx
	return nil
}

/*
================ testMongoCollection ==================
Extra methods can be added to the DBCollection interface from:
	https://github.com/mongodb/mongo-go-driver/blob/master/mongo/collection.go
as needed
See
	https://github.com/mongodb/mongo-go-driver/blob/947cf7eb5052024ab6c4ef3593d2cfb68f19e89c/x/bsonx/document.go#L48
To expand bson.Doc functionality
*/

// testMongoCollection
type testMongoCollection struct {
	name string
	ctx  context.Context
	docs []dbModel
}

// newTestMongoCollection
func newTestMongoCollection(name string) (*testMongoCollection, error) {
	if name == "" {
		return nil, errors.New("invalid test collection name")
	}
	testUserCollection := &testMongoCollection{name: name, docs: []dbModel{}}
	return testUserCollection, nil
}

// unmarshallBSON converts a BSON byte type back into a dbModel
func (coll *testMongoCollection) unmarshallBSON(bsonData interface{}) (dbModel, error) {

	return bsonUnmarshall(coll.name, bsonData)
}

// findById in the test collection a document by ID
func (coll *testMongoCollection) findById(findId string) (reDoc dbModel, err error) {
	for _, doc := range coll.docs {
		var docId string
		docId, err = standardizeID(doc)
		if err != nil {
			return
		}
		if docId == findId {
			reDoc = doc
			return
		}
	}
	return reDoc, errors.New("document not found in test collection: " + findId)
}

// deleteById in the test collection a document by ID
func (coll *testMongoCollection) deleteById(findId string) (reDoc dbModel, err error) {
	var dbDocs []dbModel
	del := false
	for _, doc := range coll.docs {
		var docId string
		docId, err = standardizeID(doc)
		if err != nil {
			return
		}
		if docId != findId {
			dbDocs = append(dbDocs, doc)
		} else {
			reDoc = doc
			del = true
		}
	}
	if !del {
		return reDoc, errors.New("document not found in test collection: " + findId)
	}
	coll.docs = dbDocs
	return reDoc, nil
}

// updateById a document in the test collection
func (coll *testMongoCollection) updateById(findId string, upDoc dbModel) (reDoc dbModel, err error) {
	var dbDocs []dbModel
	up := false
	for _, doc := range coll.docs {
		var docId string
		docId, err = standardizeID(doc)
		if err != nil {
			return
		}
		if docId != findId {
			dbDocs = append(dbDocs, doc)
		} else {
			reDoc = doc
			bsonData, bErr := upDoc.toDoc()
			if bErr != nil {
				return reDoc, bErr
			}
			err = reDoc.update(bsonData)
			if err != nil {
				return reDoc, err
			}
			up = true
			dbDocs = append(dbDocs, reDoc)
		}
	}
	if !up {
		return reDoc, errors.New("document not found in test collection: " + findId)
	}
	coll.docs = dbDocs
	return reDoc, nil
}

// find documents in the test collection
func (coll *testMongoCollection) find(dbDoc dbModel) (reDocs []dbModel, err error) {
	for _, doc := range coll.docs {
		if dbDoc == nil {
			reDocs = append(reDocs, doc)
		} else {
			bsonData, bErr := dbDoc.toDoc()
			if bErr != nil {
				return reDocs, bErr
			}
			match := false
			if len(bsonData) == 0 {
				match = true
			} else {
				match = doc.match(bsonData)
			}
			if match {
				reDocs = append(reDocs, doc)
			}
		}
	}
	return reDocs, nil
}

// insert documents into test collection
func (coll *testMongoCollection) insert(dbDocs []dbModel) (err error) {
	var valDocs []dbModel
	for _, dbDoc := range dbDocs {
		var docId string
		docId, err = standardizeID(dbDoc)
		if err != nil {
			return err
		}
		_, fErr := coll.findById(docId)
		if fErr != nil {
			// fmt.Println("----------------> CHECK THIS ERROR: ", fErr.Error())
			valDocs = append(valDocs, dbDoc)
		}
	}
	coll.docs = append(coll.docs, valDocs...)
	return nil
}

// delete documents from the test collection
func (coll *testMongoCollection) delete(dbDocs []dbModel) (reDocs []dbModel, err error) {
	for _, dbDoc := range dbDocs {
		var docId string
		docId, err = standardizeID(dbDoc)
		if err != nil {
			return reDocs, err
		}
		reDoc, fErr := coll.deleteById(docId)
		if fErr == nil {
			reDocs = append(reDocs, reDoc)
		}
	}
	return reDocs, nil
}

// InsertOne into test collection
func (coll *testMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	coll.ctx = ctx
	fmt.Println("\n--->INSERT ONE: ", document, opts)
	doc := document.(dbModel)
	err := coll.insert([]dbModel{doc})
	return &mongo.InsertOneResult{InsertedID: doc.getID()}, err
}

// InsertMany into test collection
func (coll *testMongoCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	coll.ctx = ctx
	if len(documents) == 0 {
		return nil, mongo.ErrEmptySlice
	}
	fmt.Println("\n--->INSERT MANY: ", documents, opts)
	var inDocs []dbModel
	for _, d := range documents {
		inDocs = append(inDocs, d.(dbModel))
	}
	err := coll.insert(inDocs)
	if err != nil {
		return nil, err
	}
	var inIds []interface{}
	for _, inDoc := range inDocs {
		inIds = append(inIds, inDoc.getID())
	}
	imResult := &mongo.InsertManyResult{InsertedIDs: inIds}
	return imResult, err
}

// DeleteMany from test collection
func (coll *testMongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	var delCount int64
	coll.ctx = ctx
	fmt.Println("\n--->DELETE MANY: ", filter, opts)
	filterDoc, err := coll.unmarshallBSON(filter)
	if err != nil {
		return nil, err
	}
	delDocs, err := coll.delete([]dbModel{filterDoc})
	delCount = int64(len(delDocs))
	return &mongo.DeleteResult{DeletedCount: delCount}, nil
}

// DeleteOne from test collection
func (coll *testMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	var delCount int64
	coll.ctx = ctx
	fmt.Println("\n--->DELETE ONE: ", filter, opts)
	filterDoc, err := coll.unmarshallBSON(filter)
	if err != nil {
		return nil, err
	}
	delDocs, err := coll.delete([]dbModel{filterDoc})
	delCount = int64(len(delDocs))
	return &mongo.DeleteResult{DeletedCount: delCount}, nil
}

// FindOneAndDelete finds a document, deletes it from the test collection, and then returns the found document
func (coll *testMongoCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	var rawResult []byte
	coll.ctx = ctx
	fmt.Println("\n--->FIND ONE AND DELETE: ", filter, opts)
	filterDoc, err := coll.unmarshallBSON(filter)
	if err == nil {
		delDocs, err := coll.delete([]dbModel{filterDoc})
		if err == nil && len(delDocs) > 0 {
			rawBson, err := delDocs[0].toDoc()
			if err == nil {
				rawResult, err = bsonMarshall(rawBson)
			}
		}
	}
	doc, err := bsonx.ReadDoc(rawResult)
	res := mongo.NewSingleResultFromDocument(doc, err, nil)
	return res
}

// UpdateOne a document in the test collection
func (coll *testMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	coll.ctx = ctx
	fmt.Println("\n--->UPDATE ONE: ", filter, update, opts)
	filterDoc, err := coll.unmarshallBSON(filter)
	if err != nil {
		return nil, err
	}
	docId, err := standardizeID(filterDoc)
	if err != nil {
		return nil, err
	}
	update, err = cleanUpdateBSON(update)
	if err != nil {
		panic(err)
		return nil, err
	}
	updateDoc, err := coll.unmarshallBSON(update)
	if err != nil {
		return nil, err
	}
	reDoc, err := coll.updateById(docId, updateDoc)
	return &mongo.UpdateResult{UpsertedID: reDoc.getID()}, err
}

// UpdateByID a document using an ID as the filter
func (coll *testMongoCollection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	fmt.Println("\n--->UPDATE BY ID: ", id, update, opts)
	if id == nil {
		return nil, mongo.ErrNilValue
	}
	update, err := cleanUpdateBSON(update)
	if err != nil {
		panic(err)
		return nil, err
	}
	updateDoc, err := coll.unmarshallBSON(update)
	if err != nil {
		return nil, err
	}
	return coll.UpdateOne(ctx, bson.D{{"_id", id}}, updateDoc, opts...)
}

// Find returns a collection of documents
func (coll *testMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	var rawResults []byte
	coll.ctx = ctx
	fmt.Println("\n--->FIND: ", filter, opts)
	filterDoc, err := coll.unmarshallBSON(filter)
	if err != nil {
		return nil, err
	}
	reDocs, err := coll.find(filterDoc)
	cd := initTestCursorData(reDocs)
	bsonData, err := cd.toDoc()
	if err != nil {
		panic(err)
	}
	rawResults, err = bsonMarshall(bsonData)
	cur = &mongo.Cursor{Current: rawResults}
	return cur, nil
}

// FindOne returns a single test mongo document
func (coll *testMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	var rawResult []byte
	coll.ctx = ctx
	fmt.Println("\n--->FIND ONE: ", filter, opts)
	filterDoc, err := coll.unmarshallBSON(filter)
	if err == nil {
		reDocs, err := coll.find(filterDoc)
		if err == nil && len(reDocs) > 0 {
			rawBson, err := reDocs[0].toDoc()
			if err == nil {
				rawResult, err = bsonMarshall(rawBson)
				if err != nil {
					fmt.Print("CHECK THIS ERROR: ", err.Error())
				}
			}
		}
	}
	if len(rawResult) == 0 {
		err = errors.New("document not found")
		return mongo.NewSingleResultFromDocument(nil, err, nil)
	}
	doc, _ := bsonx.ReadDoc(rawResult)
	return mongo.NewSingleResultFromDocument(doc, err, nil)
}

// CountDocuments in test mongodb collection
func (coll *testMongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	var c int64
	coll.ctx = ctx
	fmt.Println("\n--->COUNT DOCUMENTS: ", filter, opts)
	filterDoc, err := coll.unmarshallBSON(filter)
	if err != nil {
		return c, err
	}
	reDocs, err := coll.find(filterDoc)
	if err != nil {
		panic(err)
	}
	c = int64(len(reDocs))
	return c, nil
}

/*
================ testMongoDatabase ==================
*/

// testMongoDatabase
type testMongoDatabase struct {
	ctx             context.Context
	name            string
	testCollections []*testMongoCollection
}

// newTestMongoDatabase
func newTestMongoDatabase(databaseName string) (*testMongoDatabase, error) {
	if databaseName == "" {
		return &testMongoDatabase{}, errors.New("invalid test database name")
	}
	var testsColls []*testMongoCollection
	testUserCollection, err := newTestMongoCollection("users")
	if err != nil {
		fmt.Println("\nCOLLECTION INIT USER ERROR: ", err.Error())
		return &testMongoDatabase{}, err
	}
	testsColls = append(testsColls, testUserCollection)
	testGroupCollection, err := newTestMongoCollection("groups")
	if err != nil {
		fmt.Println("\nCOLLECTION INIT GROUP ERROR: ", err.Error())
		return &testMongoDatabase{}, err
	}
	testsColls = append(testsColls, testGroupCollection)
	testBlacklistCollection, err := newTestMongoCollection("blacklists")
	if err != nil {
		fmt.Println("\nCOLLECTION INIT BLACKLIST ERROR: ", err.Error())
		return &testMongoDatabase{}, err
	}
	testsColls = append(testsColls, testBlacklistCollection)
	testTasksCollection, err := newTestMongoCollection("tasks")
	if err != nil {
		fmt.Println("\nCOLLECTION INIT TASK ERROR: ", err.Error())
		return &testMongoDatabase{}, err
	}
	testsColls = append(testsColls, testTasksCollection)
	return &testMongoDatabase{
		name:            databaseName,
		testCollections: testsColls,
	}, nil
}

// Collection returns a test collection from the test client
func (c *testMongoDatabase) Collection(colName string) *testMongoCollection {
	for _, tColl := range c.testCollections {
		if tColl.name == colName {
			return tColl
		}
	}
	return nil
}

/*
================ testMongoClient ==================
*/

// testMongoClient
type testMongoClient struct {
	ctx           context.Context
	connected     bool
	testDatabases []*testMongoDatabase
}

// newTestMongoClient
func newTestMongoClient(connectionURI string) (*testMongoClient, error) {
	if connectionURI == "" {
		return nil, errors.New("invalid test connection uri")
	}
	var testDBs []*testMongoDatabase
	testMongoDB, err := newTestMongoDatabase("test")
	if err != nil {
		return nil, err
	}
	testDBs = append(testDBs, testMongoDB)
	return &testMongoClient{
		testDatabases: testDBs,
	}, nil
}

// Connect to the in-memory text mongo db
func (c *testMongoClient) Connect(ctx context.Context) error {
	c.ctx = ctx
	if c.connected {
		return errors.New("test mongo client already connected")
	}
	c.connected = true
	return nil
}

// Disconnect from the in-memory text mongo db
func (c *testMongoClient) Disconnect(ctx context.Context) error {
	c.ctx = ctx
	if !c.connected {
		return errors.New("test mongo client not connected")
	}
	c.connected = false
	return nil
}

// Database returns a test database from the test client
func (c *testMongoClient) Database(dbName string) *testMongoDatabase {
	for _, tDB := range c.testDatabases {
		if tDB.name == dbName {
			return tDB
		}
	}
	return nil
}

// Ping the in-memory text mongo db
func (c *testMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	c.ctx = ctx
	fmt.Println(rp)
	return nil
}

/*
================ testDBClient ==================
*/

// testDBClient manages a database connection
type testDBClient struct {
	connectionURI string
	client        *testMongoClient
}

// InitializeNewTestClient is a function that takes a mongoUri string and outputs a connected mongo client for the app to use
func initializeNewTestClient() (*testDBClient, error) {
	newDBClient := testDBClient{connectionURI: os.Getenv("MONGO_URI")}
	var err error
	newDBClient.client, err = newTestMongoClient(newDBClient.connectionURI)
	if err != nil {
		panic(err.Error())
	}
	return &newDBClient, err
}

// Connect opens a new connection to the database
func (db *testDBClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := db.client.Connect(ctx)
	return err
}

// Close closes an open DB connection
func (db *testDBClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := db.client.Disconnect(ctx)
	return err
}

// GetBucket returns a mongo collection based on the input collection name // todo for adding GridFS testing
func (db *testDBClient) GetBucket(bucketName string) (*gridfs.Bucket, error) {
	if bucketName == "" {
		return nil, errors.New("bucketName cannot be empty")
	}
	return nil, nil
}

// GetCollection returns a mongo collection based on the input collection name
func (db *testDBClient) GetCollection(collectionName string) DBCollection {
	return db.client.Database("test").Collection(collectionName)
}

// NewDBHandler returns a new DBHandler generic interface
func (db *testDBClient) NewDBHandler(collectionName string) *DBHandler[dbModel] {
	col := db.GetCollection(collectionName)
	return &DBHandler[dbModel]{
		db:         db,
		collection: col,
	}
}

// NewUserHandler returns a new DBHandler users interface
func (db *testDBClient) NewUserHandler() *DBHandler[*userModel] {
	col := db.GetCollection("users")
	return &DBHandler[*userModel]{
		db:         db,
		collection: col,
	}
}

// NewGroupHandler returns a new DBHandler groups interface
func (db *testDBClient) NewGroupHandler() *DBHandler[*groupModel] {
	col := db.GetCollection("groups")
	return &DBHandler[*groupModel]{
		db:         db,
		collection: col,
	}
}

// NewBlacklistHandler returns a new DBHandler blacklist interface
func (db *testDBClient) NewBlacklistHandler() *DBHandler[*blacklistModel] {
	col := db.GetCollection("blacklists")
	return &DBHandler[*blacklistModel]{
		db:         db,
		collection: col,
	}
}

// NewTaskHandler returns a new DBHandler groups interface
func (db *testDBClient) NewTaskHandler() *DBHandler[*taskModel] {
	col := db.GetCollection("tasks")
	return &DBHandler[*taskModel]{
		db:         db,
		collection: col,
	}
}

// NewFileHandler returns a new DBHandler files interface
func (db *testDBClient) NewFileHandler() *DBHandler[*fileModel] {
	col := db.GetCollection("files")
	return &DBHandler[*fileModel]{
		db:         db,
		collection: col,
	}
}
