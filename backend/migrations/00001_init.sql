-- +goose Up
CREATE TABLE ev_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    make VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    battery_kwh NUMERIC(6,2),
    range_km INT,
    price_lkr_min BIGINT,
    price_lkr_max BIGINT,
    body_type VARCHAR(30),
    import_type VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE price_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id UUID REFERENCES ev_models(id) ON DELETE CASCADE,
    source VARCHAR(50),
    price_lkr BIGINT NOT NULL,
    listed_date DATE NOT NULL,
    condition VARCHAR(20),
    mileage_km INT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE charging_stations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    operator VARCHAR(50),
    location_name VARCHAR(100),
    latitude NUMERIC(9,6),
    longitude NUMERIC(9,6),
    connector_type VARCHAR(30),
    num_chargers INT DEFAULT 1,
    is_fast_charging BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE import_policy (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    effective_date DATE NOT NULL,
    duty_rate_percent NUMERIC(5,2),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_ev_models_make_model ON ev_models (make, model);
CREATE INDEX idx_ev_models_body_type ON ev_models (body_type);
CREATE INDEX idx_ev_models_import_type ON ev_models (import_type);
CREATE INDEX idx_price_listings_model_id ON price_listings (model_id);
CREATE INDEX idx_price_listings_listed_date ON price_listings (listed_date);
CREATE INDEX idx_charging_stations_location ON charging_stations (location_name);
CREATE INDEX idx_charging_stations_connector ON charging_stations (connector_type);
CREATE INDEX idx_import_policy_effective_date ON import_policy (effective_date);

-- +goose Down
DROP TABLE IF EXISTS price_listings;
DROP TABLE IF EXISTS import_policy;
DROP TABLE IF EXISTS charging_stations;
DROP TABLE IF EXISTS ev_models;
