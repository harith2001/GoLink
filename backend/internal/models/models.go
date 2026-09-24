package models

type EVModel struct {
	ID          string   `json:"id"`
	Make        string   `json:"make"`
	Model       string   `json:"model"`
	Year        int      `json:"year"`
	BatteryKWh  *float64 `json:"battery_kwh"`
	RangeKm     *int     `json:"range_km"`
	PriceLKRMin *int64   `json:"price_lkr_min"`
	PriceLKRMax *int64   `json:"price_lkr_max"`
	BodyType    *string  `json:"body_type"`
	ImportType  *string  `json:"import_type"`
}

type PriceStats struct {
	Min   *int64 `json:"min"`
	Avg   *int64 `json:"avg"`
	Max   *int64 `json:"max"`
	Count int    `json:"count"`
}

type ModelDetails struct {
	EVModel
	PriceStats PriceStats `json:"price_stats"`
}

type PricePoint struct {
	ListedDate string  `json:"listed_date"`
	PriceLKR   int64   `json:"price_lkr"`
	Source     *string `json:"source"`
	Condition  *string `json:"condition"`
	MileageKm  *int    `json:"mileage_km"`
}

type ChargingStation struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Operator       *string  `json:"operator"`
	LocationName   *string  `json:"location_name"`
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	ConnectorType  *string  `json:"connector_type"`
	NumChargers    int      `json:"num_chargers"`
	IsFastCharging bool     `json:"is_fast_charging"`
}

type ImportPolicy struct {
	ID              string   `json:"id"`
	EffectiveDate   string   `json:"effective_date"`
	DutyRatePercent *float64 `json:"duty_rate_percent"`
	Description     *string  `json:"description"`
}
