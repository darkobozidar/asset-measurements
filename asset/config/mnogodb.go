package config

import (
    "context"
    "fmt"
    "log"
    "os"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoC *mongo.Client

func ConnectToMongoDB() {
    const BASE_CONN_URL string = "mongodb://%s:%s@mongodb:%s"

    connUri := fmt.Sprintf(
        BASE_CONN_URL,
        os.Getenv("MONGODB_INITDB_ROOT_USERNAME"),
        os.Getenv("MONGODB_INITDB_ROOT_PASSWORD"),
        os.Getenv("MONGODB_CONTAINER_PORT"),
    )

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connUri))
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
