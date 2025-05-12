package main

import (
    "asset/config"
    "asset/models"
    "asset/controllers"
    "asset/routers"

    "context"

    "github.com/gin-gonic/gin"
)

// TODO read from .env where possible.
func main() {
    // PostgreSQL
    config.ConnectToPostgresDB()
    models.MigrateModels()

    // MongoDB
    config.ConnectToMongoDB()
    config.CreateTimeSeriesCollection(
        "asset_measurements", "measurements", "timestamp", "asset_id", "seconds",
    )
    defer config.MongoC.Disconnect(context.TODO())

    // RabbitMQ
    conn := config.ConnectToRabbitMQ()
    channel := config.CreateRabbitMQChannel(conn)
    queue := config.CreateRabbitMQQueue(channel, "asset-measurements")
    go config.RegisterRabbitMQConsumer(channel, queue.Name, controllers.CreateMeasurement)
    defer conn.Close()
    defer channel.Close()

    // Start server
    r := gin.Default()
    routers.RegisterRouters(r)
    r.Run("asset:8080")
}
