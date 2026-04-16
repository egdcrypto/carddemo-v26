package repository

import (
	"context"

	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/carddemo/project/src/domain/transaction/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoTransactionRepository implements repository.TransactionRepository.
type MongoTransactionRepository struct {
	col *mongo.Collection
}

// NewMongoTransactionRepository creates a new Mongo-backed repository.
func NewMongoTransactionRepository(db *mongo.Database) repository.TransactionRepository {
	return &MongoTransactionRepository{
		col: db.Collection("transactions"),
	}
}

func (r *MongoTransactionRepository) Get(id string) (*model.Transaction, error) {
	var agg model.Transaction
	filter := bson.M{"_id": id}
	err := r.col.FindOne(context.Background(), filter).Decode(&agg)
	if err != nil {
		return nil, err
	}
	return &agg, nil
}

func (r *MongoTransactionRepository) Save(aggregate *model.Transaction) error {
	filter := bson.M{"_id": aggregate.ID}
	update := bson.M{"$set": aggregate}
	opts := bson.M{"$setOnInsert": bson.M{"created_at": aggregate.CreatedAt}}

	_, err := r.col.UpdateOne(context.Background(), filter, bson.M{"$set": update, "$setOnInsert": opts}, &mongo.UpdateOptions{
		Upsert: &[]bool{true}[0],
	})
	return err
}

func (r *MongoTransactionRepository) Delete(id string) error {
	_, err := r.col.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *MongoTransactionRepository) List() ([]*model.Transaction, error) {
	cursor, err := r.col.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*model.Transaction
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
