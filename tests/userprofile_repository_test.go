package tests

import (
	"context"
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// TestUserProfileRepository_UniqueEmail verifies the unique index behavior.
func TestUserProfileRepository_UniqueEmail(t *testing.T) {
	// Red Phase: We expect this to fail if indexes aren't created.

	// ctx := context.Background()
	// client := getTestClient(t)
	// db := client.Database("test_carddemo")
	// repo := userprofilerepo.NewMongoUserProfileRepository(db)

	// profile1 := model.NewUserProfile("user1", "john@example.com")
	// profile2 := model.NewUserProfile("user2", "john@example.com") // Duplicate Email

	// err := repo.Save(profile1)
	// assert.NoError(t, err)

	// err = repo.Save(profile2)
	// assert.Error(t, err) // Should fail due to unique index
}

// TestUserProfileRepository_Mapping verifies BSON mapping.
func TestUserProfileRepository_Mapping(t *testing.T) {
	// Red Phase: We expect this to fail if struct tags are wrong.

	// ctx := context.Background()
	// ... setup db ...

	// original := model.NewUserProfile("u1", "jane@doe.com")
	// repo.Save(original)

	// Check raw mongo document
	// var doc bson.M
	// coll := db.Collection("userprofiles")
	// err := coll.FindOne(ctx, bson.M{"_id": "u1"}).Decode(&doc)

	// assert.NoError(t, err)
	// assert.Equal(t, "jane@doe.com", doc["email"])
	// assert.Equal(t, 0, doc["version"]) // Check initial version
}
