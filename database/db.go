package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"zeotap_assign1/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Rule struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	RuleID     string             `bson:"ruleid" json:"ruleid"`
	RuleString string             `bson:"rule_string" json:"rule_string"`
	AST        *model.Node        `bson:"ast" json:"ast"`
}

var Client *mongo.Client

func ConnectDB(uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}

	fmt.Printf("Connected to MongoDB at %s\n", uri)
	return client
}

func InitializeConnections() {
	Client = ConnectDB(os.Getenv("MONGO_URI")) // Connect to the database
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	db := client.Database("ast")
	return getOrCreateCollection(db, collectionName)
}

func getOrCreateCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	exist, err := db.ListCollectionNames(context.TODO(), bson.M{"name": collectionName})
	if err != nil {
		fmt.Printf("Error listing collections: %v\n", err)
		return nil
	}

	if len(exist) == 0 {
		err := db.CreateCollection(context.TODO(), collectionName)
		if err != nil {
			fmt.Printf("Error creating collection: %v\n", err)
			return nil
		}
	}

	return db.Collection(collectionName)
}
