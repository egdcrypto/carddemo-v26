package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure implementation satisfies interface
var _ repository.AccountRepository = (*MongoAccountRepository)(nil)

// MongoAccountRepository is a MongoDB implementation of the account repository.
type MongoAccountRepository struct {
	coll *mongo.Collection
}

// NewMongoAccountRepository creates a new MongoAccountRepository.
func NewMongoAccountRepository(db *mongo.Database) *MongoAccountRepository {
	coll := db.Collection("accounts")

	// Create indexes in the background upon initialization.
	// In a production app, you might handle errors more strictly or use a migration tool.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		idxModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "_id", Value: 1}}, // MongoDB _id is unique by default, but explicit if we used custom ID field
			Options: options.Index().SetUnique(true),
		}
		// Create a unique index on the 'id' field if we store it separately from _id,
		// or rely on _id. Assuming domain ID maps to _id.
		// We also need an index on email for quick lookups if that's a query pattern.

		name, err := coll.Indexes().CreateOne(ctx, idxModel)
		if err != nil {
			fmt.Printf("[MongoAccountRepository] Failed to create index: %v\n", err)
		} else {
			fmt.Printf("[MongoAccountRepository] Created index: %s\n", name)
		}
	}()

	return &MongoAccountRepository{coll: coll}
}

// Get retrieves an account by ID.
func (r *MongoAccountRepository) Get(id string) (*model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc accountDocument
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	return doc.toDomain(), nil
}

// Save creates or updates an account aggregate.
func (r *MongoAccountRepository) Save(aggregate *model.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := fromDomain(aggregate)

	// Optimistic Locking check
	filter := bson.M{
		"_id":     doc.ID,
		"version": aggregate.Version - 1, // Check current version matches expectation
	}

	update := bson.M{
		"$set": doc,
	}

	// Upsert: Create if not exists, Update if exists and version matches.
	// If version mismatch (UpdateResult.MatchedCount == 0), return error.
	// Note: For new aggregates (Version 0), we expect UpsertedCount.
	
	opts := options.Update().SetUpsert(true)

	result, err := r.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	// Check optimistic lock failure: It existed but wasn't updated (version mismatch)
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return shared.ErrConcurrencyConflict
	}

	return nil
}

// Delete removes an account by ID.
func (r *MongoAccountRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}

	return nil
}

// List returns all accounts.
func (r *MongoAccountRepository) List() ([]*model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var docs []accountDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	accounts := make([]*model.Account, len(docs))
	for i, doc := range docs {
		accounts[i] = doc.toDomain()
	}

	return accounts, nil
}

// accountDocument represents the MongoDB schema.
type accountDocument struct {
	ID     string `bson:"_id"`
	Status string `bson:"status"`
	// Adding other fields typically found in Account Aggregate if needed
	Balance int64  `bson:"balance,omitempty"`
	Version int    `bson:"version"`
}

func fromDomain(a *model.Account) accountDocument {
	return accountDocument{
		ID:     a.ID,
		Status: a.Status,
		Version: a.Version,
		// Map other fields as necessary
	}
}

func (d accountDocument) toDomain() *model.Account {
	// We use the model constructor orhydrate the struct directly.
	// Assuming model.Account has exported fields or a Hydrate method.
	// Based on typical DDD, we might reconstruct it.
	return &model.Account{
		ID:     d.ID,
		Status: d.Status,
		Version: d.Version,
	}
}
