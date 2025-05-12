package config

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoC *mongo.Client

func ConnectToMongoDB() {
    uri := "mongodb://root:example@mongodb:27017"

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    MongoC = client
}

func CreateTimeSeriesCollection(dbName, collName, timeField, metaField, granularity string) {
    db := MongoC.Database(dbName)

    // Check if collection already exists.
    collections, err := db.ListCollectionNames(context.TODO(), map[string]interface{}{"name": collName})
    if err != nil {
        return
    }
    // Collection already exists.
    if len(collections) > 0 {
        return
    }

    tsOptions := options.TimeSeries().
        SetTimeField(timeField).
        SetMetaField(metaField).
        SetGranularity(granularity)

    opts := options.CreateCollection().SetTimeSeriesOptions(tsOptions)

    db.CreateCollection(context.TODO(), collName, opts)
}
