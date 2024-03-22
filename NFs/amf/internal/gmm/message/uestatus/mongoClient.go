package uestatus

import (
	Mongoconteext "context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
)

var client *mongo.Client

type PermanentKey struct {
	EncryptionKey       float64 `bson:"encryptionKey"`
	PermanentKeyValue   string  `bson:"permanentKeyValue"`
	EncryptionAlgorithm float64 `bson:"encryptionAlgorithm"`
}

type Milenage struct {
	Op struct {
		EncryptionAlgorithm float64 `bson:"encryptionAlgorithm"`
		EncryptionKey       float64 `bson:"encryptionKey"`
		OpValue             string  `bson:"opValue"`
	} `bson:"op"`
}

type Opc struct {
	EncryptionAlgorithm float64 `bson:"encryptionAlgorithm"`
	EncryptionKey       float64 `bson:"encryptionKey"`
	OpcValue            string  `bson:"opcValue"`
}

type AuthenticationSubscription struct {
	ID                            primitive.ObjectID `bson:"_id"`
	AuthenticationMethod          string             `bson:"authenticationMethod"`
	PermanentKey                  PermanentKey       `bson:"permanentKey"`
	SequenceNumber                string             `bson:"sequenceNumber"`
	AuthenticationManagementField string             `bson:"authenticationManagementField"`
	Milenage                      Milenage           `bson:"milenage"`
	Opc                           Opc                `bson:"opc"`
	UeID                          string             `bson:"ueId"`
}

func ConnectToMongoDB() {
	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(Mongoconteext.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to verify that the connection is established
	err = client.Ping(Mongoconteext.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB!")
}

func GetMongoData(ueId string) AuthenticationSubscription {
	// Access the MongoDB client and perform operations
	// For example, you can retrieve data from a collection here
	// This function can be called from other files in the package
	// to interact with MongoDB

	// Get a handle for the "test" database
	database := client.Database("free5gc")
	collection := database.Collection("subscriptionData.authenticationData.authenticationSubscription")
	ueid := ueId
	filter := bson.M{"ueId": ueid}

	var result AuthenticationSubscription
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func DisconnectToMongoDB() {
	// Disconnect from MongoDB when done
	err := client.Disconnect(Mongoconteext.Background())
	if err != nil {
		log.Fatal("Failed to disconnect from MongoDB:", err)
	}

	fmt.Println("Disconnected from MongoDB!")
}
