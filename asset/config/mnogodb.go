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
        os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
        os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
        os.Getenv("MONGO_CONTAINER_PORT"),
    )

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connUri))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    MongoC = client
}
