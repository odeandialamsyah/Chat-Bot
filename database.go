package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB client
var MongoClient *mongo.Client
var KeywordCollection *mongo.Collection

// Connect ke MongoDB
func ConnectDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("YOUR_MONGO_CONNECTION_STRING"))
	if err != nil {
		return err
	}

	MongoClient = client
	KeywordCollection = client.Database("whatsappdb").Collection("keywords")
	return nil
}

// Ambil kata kunci dari MongoDB
func getKeywordsFromDB() ([]Keyword, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil data dari collection keywords
	cursor, err := KeywordCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var keywords []Keyword
	for cursor.Next(ctx) {
		var keyword Keyword
		if err := cursor.Decode(&keyword); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return keywords, nil
}
