package models

import (
	"asset/config"

	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetActiveAsset(assetId uint) (Asset, error) {
	var asset Asset
	result := config.DB.First(&asset, "id = ? AND is_active = true", assetId)

	return asset, result.Error
}

func GetMongoDBAssetMeasurementsCollection() *mongo.Collection {
	db := config.MongoC.Database(os.Getenv("MONGODB_MEASUREMENTS_DB"))
	return db.Collection(os.Getenv("MONGODB_MEASUREMENTS_COLLECTION"))
}
