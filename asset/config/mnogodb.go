package config

import (
	"context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Measurement struct {
    AssetID   string    `bson:"asset_id"`  // metaField
    Timestamp time.Time `bson:"timestamp"` // timeField
    Power     float64   `bson:"power"`
    SOE       float64   `bson:"soe"`
}

// TODO replace tabs with spaces.
func ConnectToMongoDB() {
	uri := "mongodb://root:example@mongodb:27017"

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatalf("Failed to disconnect MongoDB client: %v", err)
        }
    }()

    dbName := "asset_measurements"
    collName := "measurements"

    // Create time-series collection
    if err := createTimeSeriesCollection(client, dbName, collName); err != nil {
        log.Fatalf("Error creating time-series collection: %v", err)
    }

    // Insert a sample measurement
    m := Measurement{
        AssetID:   "asset-123",
        Timestamp: time.Now(),
        Power:     456.7,
        SOE:       80.5,
    }

    if err := insertMeasurement(client, dbName, collName, m); err != nil {
        log.Fatalf("Error inserting measurement: %v", err)
    }

    fmt.Println("Measurement inserted successfully!")
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

// Insert a single measurement
func insertMeasurement(client *mongo.Client, dbName, collName string, m Measurement) error {
    collection := client.Database(dbName).Collection(collName)
    _, err := collection.InsertOne(context.TODO(), m)
    return err
}