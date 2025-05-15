CREATE TABLE IF NOT EXISTS assets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_enabled BOOLEAN,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS asset_simulation_configs (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    measurement_interval INTEGER NOT NULL,
    min_power DOUBLE PRECISION NOT NULL,
    max_power DOUBLE PRECISION NOT NULL,
    max_power_step DOUBLE PRECISION NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

-- Assets
INSERT INTO assets (name, description, type, is_enabled)
SELECT 'Battery A', 'High-efficiency lithium battery', 'battery', true
WHERE NOT EXISTS (SELECT 1 FROM assets WHERE name = 'Battery A');

INSERT INTO assets (name, description, type, is_enabled)
SELECT 'Solar Panel X', 'Rooftop panel', 'solar', true
WHERE NOT EXISTS (SELECT 1 FROM assets WHERE name = 'Solar Panel X');

INSERT INTO assets (name, description, type, is_enabled)
SELECT 'Wind Turbine 3', 'Offshore wind turbine', 'wind', false
WHERE NOT EXISTS (SELECT 1 FROM assets WHERE name = 'Wind Turbine 3');

INSERT INTO assets (name, description, type, is_enabled)
SELECT 'Battery B', 'Portable backup battery', 'battery', true
WHERE NOT EXISTS (SELECT 1 FROM assets WHERE name = 'Battery B');

INSERT INTO assets (name, description, type, is_enabled)
SELECT 'Solar Farm Y', 'Large-scale solar farm installation', 'solar', true
WHERE NOT EXISTS (SELECT 1 FROM assets WHERE name = 'Solar Farm Y');

INSERT INTO assets (name, description, type, is_enabled)
SELECT 'Wind Mill Z', 'Small wind mill in rural area', 'wind', true
WHERE NOT EXISTS (SELECT 1 FROM assets WHERE name = 'Wind Mill Z');

INSERT INTO assets (name, description, type, is_enabled)
SELECT 'Hybrid Power Box', 'Solar and battery hybrid device', 'hybrid', true
WHERE NOT EXISTS (SELECT 1 FROM assets WHERE name = 'Hybrid Power Box');

-- -- Simulation configs
INSERT INTO asset_simulation_configs (asset_id, type, measurement_interval, min_power, max_power, max_power_step)
SELECT id, 'battery', 10, -500.0, 1000.0, 50.0
FROM assets WHERE name = 'Battery A'
AND NOT EXISTS (SELECT 1 FROM asset_simulation_configs WHERE asset_id = assets.id);

INSERT INTO asset_simulation_configs (asset_id, type, measurement_interval, min_power, max_power, max_power_step)
SELECT id, 'solar', 15, -100.0, 500.0, 0.0
FROM assets WHERE name = 'Solar Panel X'
AND NOT EXISTS (SELECT 1 FROM asset_simulation_configs WHERE asset_id = assets.id);

INSERT INTO asset_simulation_configs (asset_id, type, measurement_interval, min_power, max_power, max_power_step)
SELECT id, 'wind', 20, -200.0, 1500.0, 100.0
FROM assets WHERE name = 'Wind Turbine 3'
AND NOT EXISTS (SELECT 1 FROM asset_simulation_configs WHERE asset_id = assets.id);

INSERT INTO asset_simulation_configs (asset_id, type, measurement_interval, min_power, max_power, max_power_step)
SELECT id, 'battery', 12, 150.0, 800.0, 30.0
FROM assets WHERE name = 'Battery B'
AND NOT EXISTS (SELECT 1 FROM asset_simulation_configs WHERE asset_id = assets.id);

INSERT INTO asset_simulation_configs (asset_id, type, measurement_interval, min_power, max_power, max_power_step)
SELECT id, 'solar', 20, -500.0, 300.0, 0.0
FROM assets WHERE name = 'Solar Farm Y'
AND NOT EXISTS (SELECT 1 FROM asset_simulation_configs WHERE asset_id = assets.id);

INSERT INTO asset_simulation_configs (asset_id, type, measurement_interval, min_power, max_power, max_power_step)
SELECT id, 'wind', 30, -2000.0, 1000.0, 75.0
FROM assets WHERE name = 'Wind Mill Z'
AND NOT EXISTS (SELECT 1 FROM asset_simulation_configs WHERE asset_id = assets.id);

INSERT INTO asset_simulation_configs (asset_id, type, measurement_interval, min_power, max_power, max_power_step)
SELECT id, 'hybrid', 15, -800.0, 500.0, 25.0
FROM assets WHERE name = 'Hybrid Power Box'
AND NOT EXISTS (SELECT 1 FROM asset_simulation_configs WHERE asset_id = assets.id);
