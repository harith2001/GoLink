package tools

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golink/internal/models"
)

type detailsIn struct {
	Make  string `json:"make" jsonschema:"vehicle make, e.g. BYD"`
	Model string `json:"model" jsonschema:"vehicle model, e.g. Atto 3"`
}

func registerDetails(s *mcp.Server, pool *pgxpool.Pool) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_model_details",
		Description: "Full spec for a make and model, plus min/avg/max from price_listings.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in detailsIn) (*mcp.CallToolResult, []models.ModelDetails, error) {
		rows, err := pool.Query(ctx, `
			SELECT `+modelCols+`
			FROM ev_models
			WHERE lower(make) = lower($1) AND lower(model) = lower($2)
			ORDER BY year DESC`, in.Make, in.Model)
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()
		var out []models.ModelDetails
		for rows.Next() {
			m, err := scanModel(rows)
			if err != nil {
				return nil, nil, err
			}
			d := models.ModelDetails{EVModel: m}
			if err := pool.QueryRow(ctx, `
				SELECT min(price_lkr), avg(price_lkr)::bigint, max(price_lkr), count(*)
				FROM price_listings WHERE model_id = $1`, m.ID).
				Scan(&d.PriceStats.Min, &d.PriceStats.Avg, &d.PriceStats.Max, &d.PriceStats.Count); err != nil {
				return nil, nil, err
			}
			out = append(out, d)
		}
		if err := rows.Err(); err != nil {
			return nil, nil, err
		}
		if len(out) == 0 {
			return nil, nil, fmt.Errorf("no model %s %s", in.Make, in.Model)
		}
		return nil, out, nil
	})
}
