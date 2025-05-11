package simulation

import (
	"math"
    "math/rand"
    "time"
	"log"

	"simulator/config"
	"simulator/models"
)

type AssetMeasurement struct {
    // For some reason the AssetId is not sent correctly through RabbitMQ if it
    // is named `asset_id`, but it works with `asset-id`?
    AssetID   uint      `json:"asset-id"`
    Timestamp time.Time `json:"timestamp"`
    Power     float64   `json:"power"`
    SOE       float64   `json:"soe"`
}

func GetAssetSimulationConfigs() {
	var assetSimulatorConfig []models.AssetSimulationConfig

	query := config.DB.First(&assetSimulatorConfig)
	if err := query.Error; err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
}

// func startSimulation(config AssetSimulationConfig, producer MessageProducer) {
func StartSimulation() {
	var assetConfigs []models.AssetSimulationConfig

	query := config.DB.
        // Where("isActive = true").
        Find(&assetConfigs)

	if err := query.Error; err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

    for _, assetConfig := range assetConfigs {
        // go startSimulationForConfig(assetConfig)
        log.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX %+v", assetConfig)
        go startSimulationForConfig(assetConfig)
    }
}

func startSimulationForConfig(assetConfig models.AssetSimulationConfig) {
    ticker := time.NewTicker(time.Duration(assetConfig.MeasurementInterval) * time.Second)
    defer ticker.Stop()

    // Initial state hardcoded to 50%. TODO think of something better.

    currentSOE := 50.0
    var lastPower float64 = assetConfig.MinPower

    for range ticker.C {
        currentPower := generateRandomPower(assetConfig, lastPower)
        currentSOE = updateSOE(currentSOE, currentPower, assetConfig.MeasurementInterval)
        lastPower = currentPower

        log.Printf("Power: %v, SOE: %v", currentPower, currentSOE)

        measurement := AssetMeasurement{
            AssetID:   assetConfig.ID,
            Timestamp: time.Now().UTC(),
            Power:     currentPower,
            SOE:       lastPower,
        }

        config.PublishToQueue(measurement)
    }
}

func generateRandomPower(assetConfig models.AssetSimulationConfig, lastPower float64) float64 {
    step := rand.Float64() * (assetConfig.MaxPowerStep)
    if assetConfig.MaxPowerStep <= 0 {
        step = rand.Float64() * (assetConfig.MaxPower - assetConfig.MinPower)
    }

    direction := 1.0
    if rand.Intn(2) == 0 {
        direction = -1.0
    }

    newPower := lastPower + direction*step
    newPower = math.Max(assetConfig.MinPower, math.Min(assetConfig.MaxPower, newPower))
    lastPower = newPower

    return newPower
}

func updateSOE(currentSOE, power float64, intervalSeconds int) float64 {
    // Assume full charge is 100%, min is 0%
    // Power positive = charging, negative = discharging
    delta := (power / 1000.0) * (float64(intervalSeconds) / 3600.0) * 10 // 10% per kWh/hour
    newSOE := currentSOE + delta
    return math.Max(0, math.Min(100, newSOE))
}
