package tools

import (
	"github.com/jackc/pgx/v5"
	"golink/internal/models"
)

func scanModel(row pgx.Row) (models.EVModel, error) {
	var m models.EVModel
	err := row.Scan(&m.ID, &m.Make, &m.Model, &m.Year, &m.BatteryKWh, &m.RangeKm, &m.PriceLKRMin, &m.PriceLKRMax, &m.BodyType, &m.ImportType)
	return m, err
}

const modelCols = `id::text, make, model, year, battery_kwh::float8, range_km, price_lkr_min, price_lkr_max, body_type, import_type`
