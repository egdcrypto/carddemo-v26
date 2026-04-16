package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client wraps the mongo.Client to provide application-specific connection management.
type Client struct {
	*mongo.Client
}

// NewClient creates a new MongoDB client connection.
func NewClient(uri string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB.")
	return &Client{Client: client}, nil
}

// Collection returns a reference to a collection in the specified database.
func (c *Client) Collection(dbName, collName string) *mongo.Collection {
	return c.Database(dbName).Collection(collName)
}
