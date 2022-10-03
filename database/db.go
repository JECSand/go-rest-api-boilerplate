package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
	"time"
)

// dbModel is an abstraction of the db model types
type dbModel interface {
	toDoc() (doc bson.D, err error)
	bsonFilter() (doc bson.D, err error)
	bsonUpdate() (doc bson.D, err error)
	bsonLoad(doc bson.D) (err error)
	addTimeStamps(newRecord bool)
	addObjectID()
	postProcess() (err error)
	getID() (id interface{})
	update(doc interface{}) (err error)
	match(doc interface{}) bool
}

// DBClient is an abstraction of the dbClient and testDBClient types
type DBClient interface {
	Connect() error
	Close() error
	GetBucket(bucketName string) (*gridfs.Bucket, error)
	GetCollection(collectionName string) DBCollection
	NewDBHandler(collectionName string) *DBHandler[dbModel]
	NewUserHandler() *DBHandler[*userModel]
	NewGroupHandler() *DBHandler[*groupModel]
	NewBlacklistHandler() *DBHandler[*blacklistModel]
	NewTaskHandler() *DBHandler[*taskModel]
	NewFileHandler() *DBHandler[*fileModel]
}

// DBCursor is an abstraction of the dbClient and testDBClient types
type DBCursor interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
}

// checkCursorENV returns a DBCursor based on the ENV
func checkCursorENV(cur *mongo.Cursor) DBCursor {
	if os.Getenv("ENV") == "test" {
		return newTestMongoCursor(cur)
	}
	return cur
}

// DBCollection is an abstraction of the dbClient and testDBClient types
type DBCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// DBClient manages a database connection
type dbClient struct {
	connectionURI string
	client        *mongo.Client
}

// InitializeNewClient returns an initialized DBClient based on the ENV
func InitializeNewClient() (DBClient, error) {
	if os.Getenv("ENV") == "test" {
		return initializeNewTestClient()
	}
	return initializeNewClient()
}

// InitializeNewClient is a function that takes a mongoUri string and outputs a connected mongo client for the app to use
func initializeNewClient() (*dbClient, error) {
	newDBClient := dbClient{connectionURI: os.Getenv("MONGO_URI")}
	var err error
	newDBClient.client, err = mongo.NewClient(options.Client().ApplyURI(newDBClient.connectionURI))
	return &newDBClient, err
}

// Connect opens a new connection to the database
func (db *dbClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.client.Connect(ctx)
}

// Close closes an open DB connection
func (db *dbClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.client.Disconnect(ctx)
}

// GetBucket returns a mongo collection based on the input collection name
func (db *dbClient) GetBucket(bucketName string) (*gridfs.Bucket, error) {
	bucketOpts := options.GridFSBucket()
	bucketOpts.SetName(bucketName)
	bucket, err := gridfs.NewBucket(db.client.Database(os.Getenv("DATABASE")), bucketOpts)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// GetCollection returns a mongo collection based on the input collection name
func (db *dbClient) GetCollection(collectionName string) DBCollection {
	return db.client.Database(os.Getenv("DATABASE")).Collection(collectionName)
}

// NewDBHandler returns a new DBHandler generic interface
func (db *dbClient) NewDBHandler(collectionName string) *DBHandler[dbModel] {
	col := db.GetCollection(collectionName)
	return &DBHandler[dbModel]{
		db:         db,
		collection: col,
	}
}

// NewUserHandler returns a new DBHandler users interface
func (db *dbClient) NewUserHandler() *DBHandler[*userModel] {
	col := db.GetCollection("users")
	return &DBHandler[*userModel]{
		db:         db,
		collection: col,
	}
}

// NewGroupHandler returns a new DBHandler groups interface
func (db *dbClient) NewGroupHandler() *DBHandler[*groupModel] {
	col := db.GetCollection("groups")
	return &DBHandler[*groupModel]{
		db:         db,
		collection: col,
	}
}

// NewBlacklistHandler returns a new DBHandler blacklist interface
func (db *dbClient) NewBlacklistHandler() *DBHandler[*blacklistModel] {
	col := db.GetCollection("blacklists")
	return &DBHandler[*blacklistModel]{
		db:         db,
		collection: col,
	}
}

// NewTaskHandler returns a new DBHandler task interface
func (db *dbClient) NewTaskHandler() *DBHandler[*taskModel] {
	col := db.GetCollection("tasks")
	return &DBHandler[*taskModel]{
		db:         db,
		collection: col,
	}
}

// NewFileHandler returns a new DBHandler files interface
func (db *dbClient) NewFileHandler() *DBHandler[*fileModel] {
	col := db.GetCollection("files")
	return &DBHandler[*fileModel]{
		db:         db,
		collection: col,
	}
}

// DBHandler is a Generic type struct for organizing dbModel methods
type DBHandler[T dbModel] struct {
	db         DBClient
	collection DBCollection
}

// FindOne is used to get a dbModel from the db with custom filter
func (h *DBHandler[T]) FindOne(filter T) (T, error) {
	var m T
	f, err := filter.bsonFilter()
	if err != nil {
		return filter, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = h.collection.FindOne(ctx, f).Decode(&m)
	if err != nil {
		return filter, err
	}
	return m, nil
}

// FindOneAsync is used to get a dbModel from the db with custom filter
func (h *DBHandler[T]) FindOneAsync(tCh chan T, eCh chan error, filter T, wg *sync.WaitGroup) {
	defer wg.Done()
	t, err := h.FindOne(filter)
	tCh <- t
	eCh <- err
}

// FindMany is used to get a slice of dbModels from the db with custom filter
func (h *DBHandler[T]) FindMany(filter T) ([]T, error) {
	var m []T
	f, err := filter.bsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var cur *mongo.Cursor
	if len(f) > 0 {
		cur, err = h.collection.Find(ctx, f)
	} else {
		cur, err = h.collection.Find(ctx, bson.M{})
	}
	if err != nil {
		return m, err
	}
	cursor := checkCursorENV(cur)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var md T
		cursor.Decode(&md)
		err = md.postProcess()
		if err != nil {
			return m, err
		}
		m = append(m, md)
	}
	return m, nil
}

// UpdateOne Function to update a dbModel from datasource with custom filter and update model
func (h *DBHandler[T]) UpdateOne(filter T, m T) (T, error) {
	f, err := filter.bsonFilter()
	if err != nil {
		return m, err
	}
	m.addTimeStamps(false)
	update, err := m.bsonUpdate()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = h.collection.UpdateOne(ctx, f, update)
	if err != nil {
		return m, err
	}
	err = m.postProcess()
	return m, err
}

// InsertOne adds a new dbModel record to a collection
func (h *DBHandler[T]) InsertOne(m T) (T, error) {
	m.addTimeStamps(true)
	m.addObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := h.collection.InsertOne(ctx, m)
	if err != nil {
		return m, err
	}
	err = m.postProcess()
	return m, err
}

// DeleteOne adds a new dbModel record to a collection
func (h *DBHandler[T]) DeleteOne(filter T) (T, error) { //TODO: to be replaced with "soft delete"
	var m T
	f, err := filter.bsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = h.collection.FindOneAndDelete(ctx, f).Decode(&m)
	return m, err
}

// DeleteMany adds a new dbModel record to a collection
func (h *DBHandler[T]) DeleteMany(filter T) (T, error) { //TODO: to be replaced with "soft delete"
	var m T
	f, err := filter.bsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = h.collection.DeleteMany(ctx, f)
	return filter, err
}

// newRoutine returns a new Routine for executing ASYNC DB statements
func (h *DBHandler[T]) newRoutine() *dbRoutine[T] {
	return &dbRoutine[T]{handler: h}
}
