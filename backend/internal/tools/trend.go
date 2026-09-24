package tools

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golink/internal/models"
)

type trendIn struct {
	ModelID string `json:"model_id" jsonschema:"ev_models UUID"`
	Months  *int   `json:"months,omitempty" jsonschema:"how many months back, default 12"`
}

func trendMonths(m *int) int {
	if m == nil || *m <= 0 {
		return 12
	}
	if *m > 120 {
		return 120
	}
	return *m
}

func registerTrend(s *mcp.Server, pool *pgxpool.Pool) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_price_trend",
		Description: "Listing prices for one model over the last N months.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in trendIn) (*mcp.CallToolResult, []models.PricePoint, error) {
		rows, err := pool.Query(ctx, `
			SELECT listed_date, price_lkr, source, condition, mileage_km
			FROM price_listings
			WHERE model_id = $1
			  AND listed_date >= CURRENT_DATE - make_interval(months => $2)
			ORDER BY listed_date`, in.ModelID, trendMonths(in.Months))
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()
		var out []models.PricePoint
		for rows.Next() {
			var p models.PricePoint
			var day time.Time
			if err := rows.Scan(&day, &p.PriceLKR, &p.Source, &p.Condition, &p.MileageKm); err != nil {
				return nil, nil, err
			}
			p.ListedDate = day.Format("2006-01-02")
			out = append(out, p)
		}
		return nil, out, rows.Err()
	})
}
