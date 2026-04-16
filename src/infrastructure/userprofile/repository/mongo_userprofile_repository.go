package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure implementation satisfies interface
var _ repository.UserProfileRepository = (*MongoUserProfileRepository)(nil)

// MongoUserProfileRepository is a MongoDB implementation of the userprofile repository.
type MongoUserProfileRepository struct {
	coll *mongo.Collection
}

// NewMongoUserProfileRepository creates a new MongoUserProfileRepository.
func NewMongoUserProfileRepository(db *mongo.Database) *MongoUserProfileRepository {
	coll := db.Collection("userprofiles")

	// Create unique index on email.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		idxModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		}

		name, err := coll.Indexes().CreateOne(ctx, idxModel)
		if err != nil {
			fmt.Printf("[MongoUserProfileRepository] Failed to create index: %v\n", err)
		} else {
			fmt.Printf("[MongoUserProfileRepository] Created index: %s\n", name)
		}
	}()

	return &MongoUserProfileRepository{coll: coll}
}

// Get retrieves a userprofile by ID.
func (r *MongoUserProfileRepository) Get(id string) (*model.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc userprofileDocument
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	return doc.toDomain(), nil
}

// Save creates or updates a userprofile aggregate.
func (r *MongoUserProfileRepository) Save(aggregate *model.UserProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := fromDomainProfile(aggregate)

	// Optimistic Locking
	filter := bson.M{
		"_id":     doc.ID,
		"version": aggregate.Version - 1,
	}

	update := bson.M{
		"$set": doc,
	}

	opts := options.Update().SetUpsert(true)

	result, err := r.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		// Handle duplicate key error for email explicitly if needed, or generic error
		if writeExc, ok := err.(mongo.WriteException); ok && writeExc.HasErrorCode(11000) {
			// 11000 = Duplicate Key
			return fmt.Errorf("duplicate email")
		}
		return err
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return shared.ErrConcurrencyConflict
	}

	return nil
}

// Delete removes a userprofile by ID.
func (r *MongoUserProfileRepository) Delete(id string) error {
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

// List returns all userprofiles.
func (r *MongoUserProfileRepository) List() ([]*model.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var docs []userprofileDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	profiles := make([]*model.UserProfile, len(docs))
	for i, doc := range docs {
		profiles[i] = doc.toDomain()
	}

	return profiles, nil
}

// userprofileDocument represents the MongoDB schema.
type userprofileDocument struct {
	ID       string `bson:"_id"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
	Version  int    `bson:"version"`
}

func fromDomainProfile(u *model.UserProfile) userprofileDocument {
	return userprofileDocument{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password, // In real app: ensure this is hashed logic handled before reaching repo or inside repo
		Version:  u.Version,
	}
}

func (d userprofileDocument) toDomain() *model.UserProfile {
	return &model.UserProfile{
		ID:       d.ID,
		Email:    d.Email,
		Password: d.Password,
		Version:  d.Version,
	}
}
