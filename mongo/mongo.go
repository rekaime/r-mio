package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Database interface {
	Collection(string) Collection
	Client() Client
}

type Collection interface {
	FindOne(context.Context, interface{}) *mongo.SingleResult
	Find(context.Context, interface{}) (*mongo.Cursor, error)
	InsertOne(context.Context, interface{}) (*mongo.InsertOneResult, error)
	Insert(context.Context, []interface{}) (*mongo.InsertManyResult, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Update(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(context.Context, interface{}) (int64, error)
}

type Client interface {
	Database(dbName string) Database
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
}

type Cursor interface {
}

type mongoDatabase struct {
	database *mongo.Database
}

type mongoCollection struct {
	collection *mongo.Collection
}

type mongoClient struct {
	client *mongo.Client
}

func (md *mongoDatabase) Collection(colName string) Collection {
	collection := md.database.Collection(colName)
	return &mongoCollection{collection}
}

func (md *mongoDatabase) Client() Client {
	client := md.database.Client()
	return &mongoClient{client}
}

func NewClient(uri string) (Client, error) {
	time.Local = time.UTC
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return &mongoClient{client}, err
}

func (mc *mongoClient) Database(dbName string) Database {
	database := mc.client.Database(dbName)
	return &mongoDatabase{database}
}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.client.Ping(ctx, readpref.Primary())
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.client.Disconnect(ctx)
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	return mc.collection.FindOne(ctx, filter)
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return mc.collection.Find(ctx, filter)
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return mc.collection.InsertOne(ctx, document)
}

func (mc *mongoCollection) Insert(ctx context.Context, document []interface{}) (*mongo.InsertManyResult, error) {
	return mc.collection.InsertMany(ctx, document)
}

func (mc *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.collection.UpdateOne(ctx, filter, update, opts...)
}

func (mc *mongoCollection) Update(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.collection.UpdateMany(ctx, filter, update, opts...)
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	deleteResult, err := mc.collection.DeleteOne(ctx, filter)
	return deleteResult.DeletedCount, err
}
