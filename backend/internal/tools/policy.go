package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golink/internal/models"
)

type policyIn struct {
	AsOfDate *string `json:"as_of_date,omitempty" jsonschema:"YYYY-MM-DD; omit for the latest row"`
}

func registerPolicy(s *mcp.Server, pool *pgxpool.Pool) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_import_policy",
		Description: "Duty rate in effect on a date, or the latest row when the date is omitted. Figures are curated demo values.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in policyIn) (*mcp.CallToolResult, *models.ImportPolicy, error) {
		var day *time.Time
		if in.AsOfDate != nil && *in.AsOfDate != "" {
			t, err := time.Parse("2006-01-02", *in.AsOfDate)
			if err != nil {
				return nil, nil, fmt.Errorf("as_of_date must be YYYY-MM-DD")
			}
			day = &t
		}
		var p models.ImportPolicy
		var eff time.Time
		err := pool.QueryRow(ctx, `
			SELECT id::text, effective_date, duty_rate_percent::float8, description
			FROM import_policy
			WHERE ($1::date IS NULL OR effective_date <= $1)
			ORDER BY effective_date DESC
			LIMIT 1`, day).Scan(&p.ID, &eff, &p.DutyRatePercent, &p.Description)
		if err != nil {
			return nil, nil, err
		}
		p.EffectiveDate = eff.Format("2006-01-02")
		return nil, &p, nil
	})
}
