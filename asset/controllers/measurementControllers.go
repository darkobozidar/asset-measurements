package controllers

import (
	"asset/config"

	"net/http"
	"strconv"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Measurement struct {
	// For some reason the AssetId is not sent correctly through RabbitMQ if it
    // is named `asset_id`, but it works with `asset-id`?
    AssetID   uint      `bson:"asset_id""`  // metaField
    Timestamp time.Time `bson:"timestamp"` // timeField
    Power     float64   `bson:"power"`
    SOE       float64   `bson:"soe"`
}

// TODO check to create indexes
func GetLatestMeasurement(c *gin.Context) {
    assetID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset_id"})
        return
    }

    collection := config.CLIENT.Database("asset_measurements").Collection("measurements")

    filter := bson.M{"asset_id": assetID}
    opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})

    var measurement Measurement
    err = collection.FindOne(context.TODO(), filter, opts).Decode(&measurement)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "No measurements found for this asset"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }

    c.JSON(http.StatusOK, measurement)
}

func GetMeasurementsInRange(c *gin.Context) {
    assetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset_id"})
        return
    }

    fromStr := c.Query("from")
    toStr := c.Query("to")
    order := c.DefaultQuery("order", "asc")

    from, err1 := time.Parse(time.RFC3339, fromStr)
    to, err2 := time.Parse(time.RFC3339, toStr)
    if err1 != nil || err2 != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from/to date format (expected RFC3339)"})
        return
    }

    sortOrder := 1
    if order == "desc" {
        sortOrder = -1
    }

    collection := config.CLIENT.Database("asset_measurements").Collection("measurements")
    filter := bson.M{
        "asset_id": assetID,
        "timestamp": bson.M{
            "$gte": from,
            "$lte": to,
        },
    }

    cursor, err := collection.Find(
        context.TODO(),
        filter,
        options.Find().SetSort(bson.D{{"timestamp", sortOrder}}),
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
        return
    }

    var results []Measurement
    if err := cursor.All(context.TODO(), &results); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "decoding failed"})
        return
    }

    c.JSON(http.StatusOK, results)
}

func GetAverageMeasurements(c *gin.Context) {
    assetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset_id"})
        return
    }

    fromStr := c.Query("from")
    toStr := c.Query("to")
    interval := c.DefaultQuery("interval", "hour") // default
    order := c.DefaultQuery("order", "asc")

    from, err1 := time.Parse(time.RFC3339, fromStr)
    to, err2 := time.Parse(time.RFC3339, toStr)
    if err1 != nil || err2 != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from/to date format (expected RFC3339)"})
        return
    }

    var format string
    switch interval {
    case "minute":
        format = "%Y-%m-%dT%H:%M"
    case "15min":
        // We'll bucket manually into 15-minute slots using a $function or math (see below)
        // But since $function is not available in Go driver directly, we'll round timestamp to nearest 15min in code.
        format = "" // we’ll use $dateTrunc instead (MongoDB 5.0+)
    default: // "hour"
        format = "%Y-%m-%dT%H"
    }

    sortOrder := 1
    if order == "desc" {
        sortOrder = -1
    }

    collection := config.CLIENT.Database("asset_measurements").Collection("measurements")

    var timeGroup bson.D
    if interval == "15min" {
        // MongoDB >=5.0: $dateTrunc for accurate 15-minute rounding
        timeGroup = bson.D{
            {"$dateTrunc", bson.D{
                {"date", "$timestamp"},
                {"unit", "minute"},
                {"binSize", 15},
            }},
        }
    } else {
        timeGroup = bson.D{
            {"$dateToString", bson.D{
                {"format", format},
                {"date", "$timestamp"},
            }},
        }
    }

    pipeline := mongo.Pipeline{
        bson.D{{"$match", bson.D{
            {"asset_id", assetID},
            {"timestamp", bson.D{{"$gte", from}, {"$lte", to}}},
        }}},
        bson.D{{"$group", bson.D{
            {"_id", timeGroup},
            {"avg_power", bson.D{{"$avg", "$power"}}},
            {"avg_soe", bson.D{{"$avg", "$soe"}}},
        }}},
        bson.D{{"$sort", bson.D{{"_id", sortOrder}}}},
    }

    cursor, err := collection.Aggregate(context.TODO(), pipeline)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "aggregation failed"})
        return
    }

    var results []bson.M
    if err := cursor.All(context.TODO(), &results); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "decoding failed"})
        return
    }

    c.JSON(http.StatusOK, results)
}
