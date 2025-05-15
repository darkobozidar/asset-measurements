package simulation

import (
	"math"
    "math/rand"
    "time"
    "sync"
    "context"

	"simulator/models"
)

type assetMeasurement struct {
    AssetID   uint      `json:"asset_id"`
    Timestamp time.Time `json:"timestamp"`
    Power     float64   `json:"power"`
    SOE       float64   `json:"soe"`
}

// Map currently used only for adding simulations, but not for stopping them.
// This is groundwork for the assignment expansion in case we wanted to change
// the existing simulations (eg. add new, edit / delete existing).
type SimulationManager struct {
    // Using sync.Map to avoid panics for concurrent read / write to a map.
	Simulations sync.Map
}

func (sm *SimulationManager) StartSimulation(simulationHandler func(obj any)) {
	assetSimulationConfigs := models.GetActiveAssetSimulationConfigs()

    for _, assetConfig := range assetSimulationConfigs {
        sm.startSimulationForAsset(assetConfig, simulationHandler)
    }
}

func (sm *SimulationManager) startSimulationForAsset(assetConfig models.AssetSimulationConfig, simulationHandler func(obj any)) {
    if _, exists := sm.Simulations.Load(assetConfig.AssetID); exists {
        return
    }

    _, cancel := context.WithCancel(context.Background())
    previousPower := (assetConfig.MaxPower - assetConfig.MinPower) / 2
    previousSOE := 50.0

    go func() {
        ticker := time.NewTicker(time.Duration(assetConfig.MeasurementInterval) * time.Second)
        defer ticker.Stop()

        for range ticker.C {
            currentPower := generatePower(previousPower, assetConfig)
            currentSOE := generateSOE(currentPower, previousSOE, float64(assetConfig.MeasurementInterval))

            measurement := assetMeasurement{
                AssetID:   assetConfig.ID,
                Timestamp: time.Now().UTC(),
                Power:     currentPower,
                SOE:       currentSOE,
            }

            previousPower = currentPower
            previousSOE = currentSOE

            simulationHandler(measurement)
        }
    }()

    sm.Simulations.Store(assetConfig.AssetID, cancel)
}

func generatePower(previousPower float64, assetConfig models.AssetSimulationConfig) float64 {
    step := rand.Float64() * (assetConfig.MaxPowerStep)
    if assetConfig.MaxPowerStep <= 0 {
        step = rand.Float64() * (assetConfig.MaxPower - assetConfig.MinPower)
    }

    direction := 1.0
    if rand.Intn(2) == 0 {
        direction = -1.0
    }

    currentPower := previousPower + direction * step
    currentPower = math.Max(assetConfig.MinPower, math.Min(assetConfig.MaxPower, currentPower))

    return currentPower
}

func generateSOE(currentPower, previousSOE, measurementInterval float64, ) float64 {
    // For simulation simplicity, lets assume that power change always affects the SEO for
    // the same constant `POWER_TO_SOE_RATIO`.
    const POWER_TO_SOE_RATIO float64 = 0.0001
    delta := currentPower * measurementInterval * POWER_TO_SOE_RATIO

    currentSOE := previousSOE + delta
    currentSOE = math.Max(0, math.Min(100, currentSOE))

    return currentSOE
}
