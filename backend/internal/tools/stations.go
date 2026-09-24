package tools

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golink/internal/models"
)

type stationsIn struct {
	NearLocation     *string `json:"near_location,omitempty" jsonschema:"substring of location_name, e.g. Colombo"`
	ConnectorType    *string `json:"connector_type,omitempty" jsonschema:"CCS2, CHAdeMO, or Type2"`
	FastChargingOnly *bool   `json:"fast_charging_only,omitempty" jsonschema:"when true, only DC fast chargers"`
}

func registerStations(s *mcp.Server, pool *pgxpool.Pool) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "find_charging_stations",
		Description: "Charging stations by place name, connector, and fast-charging flag.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in stationsIn) (*mcp.CallToolResult, []models.ChargingStation, error) {
		rows, err := pool.Query(ctx, `
			SELECT id::text, name, operator, location_name,
			       latitude::float8, longitude::float8, connector_type,
			       num_chargers, is_fast_charging
			FROM charging_stations
			WHERE ($1::text IS NULL OR location_name ILIKE '%' || $1 || '%')
			  AND ($2::text IS NULL OR lower(connector_type) = lower($2))
			  AND ($3::bool IS NOT TRUE OR is_fast_charging)
			ORDER BY location_name, name
			LIMIT 50`, in.NearLocation, in.ConnectorType, in.FastChargingOnly)
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()
		var out []models.ChargingStation
		for rows.Next() {
			var st models.ChargingStation
			if err := rows.Scan(&st.ID, &st.Name, &st.Operator, &st.LocationName, &st.Latitude, &st.Longitude, &st.ConnectorType, &st.NumChargers, &st.IsFastCharging); err != nil {
				return nil, nil, err
			}
			out = append(out, st)
		}
		return nil, out, rows.Err()
	})
}
