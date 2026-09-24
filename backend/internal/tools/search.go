package tools

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golink/internal/models"
)

type searchIn struct {
	MaxBudgetLKR *int64  `json:"max_budget_lkr,omitempty" jsonschema:"highest acceptable price_lkr_min"`
	MinRangeKm   *int    `json:"min_range_km,omitempty" jsonschema:"minimum range in km"`
	BodyType     *string `json:"body_type,omitempty" jsonschema:"hatchback, sedan, suv, or crossover"`
	ImportType   *string `json:"import_type,omitempty" jsonschema:"brand_new or used"`
}

func registerSearch(s *mcp.Server, pool *pgxpool.Pool) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_models",
		Description: "Find Sri Lanka EV models by budget, range, body type, and import type.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in searchIn) (*mcp.CallToolResult, []models.EVModel, error) {
		rows, err := pool.Query(ctx, `
			SELECT `+modelCols+`
			FROM ev_models
			WHERE ($1::bigint IS NULL OR price_lkr_min <= $1)
			  AND ($2::int IS NULL OR range_km >= $2)
			  AND ($3::text IS NULL OR lower(body_type) = lower($3))
			  AND ($4::text IS NULL OR lower(import_type) = lower($4))
			ORDER BY price_lkr_min NULLS LAST, make, model
			LIMIT 50`, in.MaxBudgetLKR, in.MinRangeKm, in.BodyType, in.ImportType)
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()
		var out []models.EVModel
		for rows.Next() {
			m, err := scanModel(rows)
			if err != nil {
				return nil, nil, err
			}
			out = append(out, m)
		}
		return nil, out, rows.Err()
	})
}
