package mocks

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockCollection is a mock for mongo.Collection interface.
// Note: mongo.Collection is an interface, but its methods return concrete types
// (Cursor, SingleResult) which are structs, making pure interface mocking hard
// without generics or helper wrappers. 
// In a real scenario, we wrap the driver or use a test container.
// For this task, we rely on the prompt's instruction to use mock adapters.

// We will define a simplified interface that our Repository uses, and mock that.

// Collection is a subset of mongo.Collection methods used by our repos.
type Collection interface {
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Indexes() IndexView
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

type IndexView interface {
	CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
	CreateMany(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error)
}

// MockCollection implementation for testing.
type MockCollection struct {
	InsertOneFunc func(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	FindOneFunc   func(ctx context.Context, filter interface{}) *mongo.SingleResult
	UpdateOneFunc func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOneFunc func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	IndexesFunc   func() IndexView
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if m.InsertOneFunc != nil {
		return m.InsertOneFunc(ctx, document)
	}
	return &mongo.InsertOneResult{}, nil
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	if m.FindOneFunc != nil {
		return m.FindOneFunc(ctx, filter)
	}
	// Return an empty SingleResult or decode error mock
	return nil
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if m.UpdateOneFunc != nil {
		return m.UpdateOneFunc(ctx, filter, update, opts...)
	}
	return &mongo.UpdateResult{}, nil
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	if m.DeleteOneFunc != nil {
		return m.DeleteOneFunc(ctx, filter)
	}
	return &mongo.DeleteResult{}, nil
}

func (m *MockCollection) Indexes() IndexView {
	if m.IndexesFunc != nil {
		return m.IndexesFunc()
	}
	return &MockIndexView{}
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return nil, nil
}

// MockIndexView
type MockIndexView struct {
	CreateOneFunc  func(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
	CreateManyFunc func(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error)
}

func (m *MockIndexView) CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	if m.CreateOneFunc != nil {
		return m.CreateOneFunc(ctx, model, opts...)
	}
	return "idx_123", nil
}

func (m *MockIndexView) CreateMany(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	if m.CreateManyFunc != nil {
		return m.CreateManyFunc(ctx, models, opts...)
	}
	return []string{"idx_123"}, nil
}
