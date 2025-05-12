package models

import (
	"simulator/utils"

	"simulator/config"
)

func GetActiveAssetSimulationConfigs() []AssetSimulationConfig {
	var assetSimulationConfigs []AssetSimulationConfig

	query := config.DB.
        Where("is_active = true").
        Find(&assetSimulationConfigs)
	utils.FailOnError(query.Error, "Error while reading active AssetSimulationConfig records.")

	return assetSimulationConfigs
}
