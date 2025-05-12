package config

import (
	"context"
    "fmt"
    "log"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoC *mongo.Client

// TODO replace tabs with spaces.
func ConnectToMongoDB() {
	uri := "mongodb://root:example@mongodb:27017"

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // TODO clean this up at the end.
    // defer func() {
    //     if err = client.Disconnect(context.TODO()); err != nil {
    //         log.Fatalf("Failed to disconnect MongoDB client: %v", err)
    //     }
    // }()

    dbName := "asset_measurements"
    collName := "measurements"

    // Create time-series collection
    if err := createTimeSeriesCollection(client, dbName, collName); err != nil {
        log.Fatalf("Error creating time-series collection: %v", err)
    }

	MongoC = client
}

// Create the time-series collection
func createTimeSeriesCollection(client *mongo.Client, dbName, collName string) error {
    db := client.Database(dbName)

    // Check if collection already exists
    collections, err := db.ListCollectionNames(context.TODO(), map[string]interface{}{"name": collName})
    if err != nil {
        return err
    }
    if len(collections) > 0 {
        fmt.Println("Collection already exists. Skipping creation.")
        return nil
    }

    tsOptions := options.TimeSeries().
        SetTimeField("timestamp").
        SetMetaField("asset_id").
        SetGranularity("seconds")

    opts := options.CreateCollection().SetTimeSeriesOptions(tsOptions)

    return db.CreateCollection(context.TODO(), collName, opts)
}
