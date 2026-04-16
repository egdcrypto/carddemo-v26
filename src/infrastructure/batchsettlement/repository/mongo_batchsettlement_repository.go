package repository

import (
	"context"

	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/batchsettlement/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoBatchSettlementRepository implements repository.BatchSettlementRepository.
type MongoBatchSettlementRepository struct {
	col *mongo.Collection
}

// NewMongoBatchSettlementRepository creates a new Mongo-backed repository.
func NewMongoBatchSettlementRepository(db *mongo.Database) repository.BatchSettlementRepository {
	return &MongoBatchSettlementRepository{
		col: db.Collection("settlements"),
	}
}

func (r *MongoBatchSettlementRepository) Get(id string) (*model.BatchSettlement, error) {
	var agg model.BatchSettlement
	filter := bson.M{"_id": id}
	err := r.col.FindOne(context.Background(), filter).Decode(&agg)
	if err != nil {
		return nil, err
	}
	return &agg, nil
}

func (r *MongoBatchSettlementRepository) Save(aggregate *model.BatchSettlement) error {
	filter := bson.M{"_id": aggregate.ID}
	update := bson.M{"$set": aggregate}
	opts := bson.M{"$setOnInsert": bson.M{"created_at": aggregate.CreatedAt}}

	_, err := r.col.UpdateOne(context.Background(), filter, bson.M{"$set": update, "$setOnInsert": opts}, &mongo.UpdateOptions{
		Upsert: &[]bool{true}[0],
	})
	return err
}

func (r *MongoBatchSettlementRepository) Delete(id string) error {
	_, err := r.col.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *MongoBatchSettlementRepository) List() ([]*model.BatchSettlement, error) {
	cursor, err := r.col.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*model.BatchSettlement
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
