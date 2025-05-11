package models

import (
	"asset/config"
)

type Asset struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null" json:"name"`
    Description string    `gorm:"type:text" json:"description"`
    Type        string    `gorm:"type:varchar(50);not null" json:"type"`
    Enabled     bool      `gorm:"default:true" json:"enabled"`
    // CreatedAt   time.Time `json:"created_at"`  // TODO
    // UpdatedAt   time.Time `json:"updated_at"`
}

func MigrateModels() {
	config.DB.AutoMigrate(&Asset{})
}
