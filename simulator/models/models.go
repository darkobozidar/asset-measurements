package models

type AssetSimulationConfig struct {
    ID                  uint      `gorm:"primaryKey" json:"id"`
    AssetID             uint      `gorm:"not null" json:"asset_id"`
    Type                string    `gorm:"type:varchar(50);not null" json:"type"`
    MeasurementInterval int       `gorm:"not null" json:"measurement_interval"`
    MinPower            float64   `gorm:"not null" json:"min_power"`
    MaxPower            float64   `gorm:"not null" json:"max_power"`
    MaxPowerStep        float64   `gorm:"not null" json:"max_power_step"`
    IsActive            bool      `gorm:"default:true" json:"-"`
}
