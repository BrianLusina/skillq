package mongodb

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoDBClient is a structure for handling mongo db connections. This contains the underlying mongo client, database and collection
// each instance might have a different client database and collection each with different settings depending on the configuration provided
type mongoDBClient[T any] struct {
	mongoClient *mongo.Client
	database    *mongo.Database
	collection  *mongo.Collection
}

// New creates a new mongo DB client
func New[T any](config MongoDBConfig) (MongoDBClient[T], error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Client.User, config.Client.Password, config.Client.Host, config.Client.Port)
	clientOptions := options.Client().ApplyURI(uri)

	clientOptions.Hosts = []string{config.Client.Host}
	clientOptions.SetRetryWrites(config.Client.RetryWrites)

	dbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to connect to mongo db: %v", err))
		return nil, errors.Wrapf(err, "failed to connect to mongo DB")
	}

	// TODO: set database options if provided
	dbOptions := options.Database()

	db := dbClient.Database(config.DBConfig.DatabaseName, dbOptions)
	collection := db.Collection(config.DBConfig.CollectionName)

	slog.Info("connected to mongo db")

	return &mongoDBClient[T]{
		mongoClient: dbClient,
		database:    db,
		collection:  collection,
	}, nil
}

// Insert inserts a given model to the database's collection & returns the ID /error if any
func (client *mongoDBClient[T]) Insert(ctx context.Context, model T) (primitive.ObjectID, error) {
	result, err := client.collection.InsertOne(ctx, model)
	if err != nil {
		slog.ErrorContext(ctx, "failed to insert item: %v with err: %v", model, err)
		return primitive.ObjectID{}, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("failed to retrieve object id from result: %v", result)
	}

	return oid, nil
}

func (client *mongoDBClient[T]) BulkInsert(ctx context.Context, models []any) ([]primitive.ObjectID, error) {
	result, err := client.collection.InsertMany(ctx, models)
	if err != nil {
		return nil, err
	}
	insertedIds := []primitive.ObjectID{}

	for _, insertedId := range result.InsertedIDs {
		oid, ok := insertedId.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve object id from result: %v", insertedId)
		}
		insertedIds = append(insertedIds, oid)
	}

	return insertedIds, nil
}

func (client *mongoDBClient[T]) Delete(ctx context.Context, keyName string, id string) error {
	filter := bson.D{{
		Key:   keyName,
		Value: id,
	}}

	result := client.collection.FindOneAndDelete(ctx, filter)

	var d bson.D
	err := result.Decode(&d)
	if err != nil {
		return err
	}

	return nil
}

func (client *mongoDBClient[T]) FindById(ctx context.Context, keyName string, id string) (T, error) {
	filter := bson.D{{
		Key:   keyName,
		Value: id,
	}}

	var result bson.D
	err := client.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return *new(T), err
	}

	document, err := bson.Marshal(&result)
	if err != nil {
		return *new(T), err
	}

	var model T
	err = bson.Unmarshal(document, model)
	if err != nil {
		return *new(T), err
	}

	return model, nil
}

func (client *mongoDBClient[T]) FindAll(ctx context.Context, filter map[string]map[string]string) ([]T, error) {
	filterValues := bson.M{}

	for key, value := range filter {
		nestedBsonMap := bson.M{}

		for nestedKey, nestedValue := range value {
			nestedBsonMap[nestedKey] = nestedValue
		}

		filterValues[key] = nestedBsonMap
	}

	cursor, err := client.collection.Find(ctx, filterValues)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := []T{}
	for cursor.Next(ctx) {
		var model T
		err := cursor.Decode(&model)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode result")
		}
		results = append(results, model)
	}

	return results, nil
}

// Disconnect disconnects from a mongo db client connection
func (client *mongoDBClient[T]) Disconnect(ctx context.Context) error {
	return client.mongoClient.Disconnect(ctx)
}
