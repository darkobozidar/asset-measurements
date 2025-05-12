package models

import (
	"asset/config"

    "time"
)

type Asset struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null" json:"name"`
    Description string    `gorm:"type:text" json:"description"`
    Type        string    `gorm:"type:varchar(50);not null" json:"type"`
    IsEnabled   bool      `json:"isEnabled"`
    IsActive    bool      `gorm:"default:true" json:"-"`
}

type AssetMeasurement struct {
    AssetID   uint      `bson:"asset_id" json:"asset_id"`
    Timestamp time.Time `bson:"timestamp"`
    Power     float64   `bson:"power"`
    SOE       float64   `bson:"soe"`
}

func MigrateModels() {
	config.DB.AutoMigrate(&Asset{})
}
