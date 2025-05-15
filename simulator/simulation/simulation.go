package simulation

import (
	"math"
    "math/rand"
    "time"

	"simulator/models"
)

type assetMeasurement struct {
    AssetID   uint      `json:"asset_id"`
    Timestamp time.Time `json:"timestamp"`
    Power     float64   `json:"power"`
    SOE       float64   `json:"soe"`
}

func StartSimulation(simulationHandler func(obj any)) {
	assetSimulationConfigs := models.GetActiveAssetSimulationConfigs()

    for _, assetConfig := range assetSimulationConfigs {
        go startSimulationForConfig(assetConfig, simulationHandler)
    }
}

func startSimulationForConfig(assetConfig models.AssetSimulationConfig, simulationHandler func(obj any)) {
    ticker := time.NewTicker(time.Duration(assetConfig.MeasurementInterval) * time.Second)
    defer ticker.Stop()

    // Initial state hardcoded to 50%.
    currentSOE := 50.0
    var lastPower float64 = assetConfig.MinPower

    for range ticker.C {
        currentPower := generateRandomPower(assetConfig, lastPower)
        currentSOE = updateSOE(currentSOE, currentPower, assetConfig.MeasurementInterval)
        lastPower = currentPower

        measurement := assetMeasurement{
            AssetID:   assetConfig.ID,
            Timestamp: time.Now().UTC(),
            Power:     currentPower,
            SOE:       lastPower,
        }

        simulationHandler(measurement)
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
