package main

import (
    "asset/config"
    "asset/controllers"
    "asset/routers"

    "context"
    "os"
    "fmt"

    "github.com/gin-gonic/gin"
)

func main() {
    // PostgreSQL
    config.ConnectToPostgresDB()
    defer config.SQLDB.Close()

    // MongoDB
    TIME_FIELD := "timestamp"
    META_FIELD := "asset_id"
    GRANULARITY := "seconds"
    config.ConnectToMongoDB()
    config.CreateTimeSeriesCollection(
        os.Getenv("MONGODB_MEASUREMENTS_DB"),
        os.Getenv("MONGODB_MEASUREMENTS_COLLECTION"),
        TIME_FIELD,
        META_FIELD,
        GRANULARITY,
    )
    defer config.MongoC.Disconnect(context.TODO())

    // RabbitMQ
    conn := config.ConnectToRabbitMQ()
    channel := config.CreateRabbitMQChannel(conn)
    queue := config.CreateRabbitMQQueue(channel, os.Getenv("RABBITMQ_MEASUREMENTS_QUEUE_NAME"))
    go config.RegisterRabbitMQConsumer(channel, queue.Name, controllers.CreateMeasurement)
    defer conn.Close()
    defer channel.Close()

    // Start server
    r := gin.Default()
    routers.RegisterRouters(r)
    r.Run(fmt.Sprintf("%s:%s", os.Getenv("ASSET_SERVICE_HOST"), os.Getenv("ASSET_SERVICE_CONTAINER_PORT")))
}
