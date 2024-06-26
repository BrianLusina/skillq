package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoDBClient defines the capabilities of this mongo database client
type MongoDBClient[T any] interface {
	// Insert adds a new data item
	Insert(ctx context.Context, model T) (primitive.ObjectID, error)

	// BulkInsert inserts a collection of data items
	BulkInsert(ctx context.Context, models []any) ([]primitive.ObjectID, error)

	// FindById retrieves the model using the given key name and the id value populating the model with the retrieved data if found
	FindById(ctx context.Context, keyName string, id string) (T, error)

	// FindAll retrieves all the items with a given filter
	FindAll(ctx context.Context, filterOptions FilterOptions) ([]T, error)

	// Update updates the model
	Update(ctx context.Context, model T, updateOptions UpdateOptions) error

	// Delete deletes a record given it's ID name and the id value
	Delete(ctx context.Context, keyName string, id string) error

	// Disconnect disconnects from the current connection
	Disconnect(context.Context) error

	// CreateIndex creates unique index for the given keys
	CreateIndex(ctx context.Context, indexParam IndexParam) (string, error)
}
