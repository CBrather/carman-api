package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

// GetDBClient returns the db, opening the connection first if not already done.
func GetDBClient(connectionString string) (*mongo.Client, error) {
	if dbClient != nil {
		return dbClient, nil
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	connectionOptions := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)

	dbClient, err := mongo.Connect(context.TODO(), connectionOptions)
	if err != nil {
		return nil, err
	}

	var pingResult bson.M
	err = dbClient.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&pingResult)
	if err != nil {
		return nil, err
	}

	return dbClient, nil
}
