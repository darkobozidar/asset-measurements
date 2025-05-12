package models

import (
	"asset/config"
)

func GetActiveAsset(assetId uint) (Asset, error) {
	var asset Asset
	result := config.DB.First(&asset, "id = ? AND is_active = true", assetId)

	return asset, result.Error
}
