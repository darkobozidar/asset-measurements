package controllers

import (
    "asset/config"
    "asset/utils"
    "asset/models"

    "encoding/json"
    "net/http"
    "strconv"
    "context"
    "time"
    "log"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// TODO check to create indexes.
func GetLatestMeasurement(c *gin.Context) {
    assetID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset_id"})
        return
    }

    collection := config.MongoC.Database("asset_measurements").Collection("measurements")

    filter := bson.M{"asset_id": assetID}
    opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})

    var measurement models.AssetMeasurement
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

    collection := config.MongoC.Database("asset_measurements").Collection("measurements")
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

    var results []models.AssetMeasurement
    if err := cursor.All(context.TODO(), &results); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "decoding failed"})
        return
    }

    c.JSON(http.StatusOK, results)
}

func GetAverageMeasurements(c *gin.Context) {
    var queryParams struct {
        From     string `form:"from" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
        To       string `form:"to" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
        GroupBy  string    `form:"groupBy" binding:"omitempty,oneof=1minute 15minute 1hour"`
        Sort     string    `form:"sort" binding:"omitempty,oneof=asc desc"`
    }

    if err := c.ShouldBindQuery(&queryParams); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    assetID, err := utils.StringToUint(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fromDateTime, err := time.Parse(time.RFC3339, queryParams.From)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date format (expected RFC3339)"})
    }
    toDateTime, err := time.Parse(time.RFC3339, queryParams.To)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date format (expected RFC3339)"})
    }

    binSize, unit := 1, "minute"
    if queryParams.GroupBy != "" {
        binSize, unit, err = utils.ExtractBinSizeAndUnit(queryParams.GroupBy)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    }

    sortOrder := 1
    if queryParams.Sort == "desc" {
        sortOrder = -1
    }

    // TODO fix this.
    collection := config.MongoC.Database("asset_measurements").Collection("measurements")

    pipeline := mongo.Pipeline{
        bson.D{
            {"$match", bson.D{
                {"asset_id", assetID},
                {"timestamp", bson.D{
                    {"$gte", fromDateTime},
                    {"$lte", toDateTime},
                }},
            }},
        },
        bson.D{
            {"$group", bson.D{
                {"_id", bson.D{
                    {"$dateTrunc", bson.D{
                        {"date", "$timestamp"},
                        {"unit", unit},
                        {"binSize", binSize},
                    }},
                }},
                {"avg_power", bson.D{{"$avg", "$power"}}},
                {"avg_soe", bson.D{{"$avg", "$soe"}}},
            }},
        },
        bson.D{
            {"$sort", bson.D{
                {"_id", sortOrder},
            }},
        },
    }

    cursor, err := collection.Aggregate(context.TODO(), pipeline)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Aggregation failed"})
        return
    }

    var results []bson.M
    if err := cursor.All(context.TODO(), &results); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Decoding failed"})
        return
    }

    c.JSON(http.StatusOK, results)
}

func CreateMeasurement(msg []byte) {
    var assetMeasurement models.AssetMeasurement
    err := json.Unmarshal(msg, &assetMeasurement)
    utils.FailOnError(err, "Failed to decode JSON")

    asset, err := models.GetActiveAsset(assetMeasurement.AssetID)
    if err != nil {
        utils.LogOnError(err, "Error while reading active Asset.")
        return
    }

    if !asset.IsEnabled {
        log.Printf("Asset %+v disabled. Not saving to DB the measurement %+v.", asset, assetMeasurement)
        return
    }

    // TODO check if assetID exists (Exactly once!)
    collection := config.MongoC.Database("asset_measurements").Collection("measurements")
    _, err = collection.InsertOne(context.TODO(), assetMeasurement)
    utils.FailOnError(err, "Error on inserting measurement");
}
