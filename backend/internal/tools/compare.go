package tools

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golink/internal/models"
)

type compareIn struct {
	ModelIDs []string `json:"model_ids" jsonschema:"two or three ev_models UUIDs"`
}

func checkCompareCount(n int) error {
	if n < 2 || n > 3 {
		return fmt.Errorf("model_ids must contain 2 or 3 ids, got %d", n)
	}
	return nil
}

func registerCompare(s *mcp.Server, pool *pgxpool.Pool) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "compare_models",
		Description: "Side-by-side specs for 2 or 3 models by id.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in compareIn) (*mcp.CallToolResult, []models.EVModel, error) {
		if err := checkCompareCount(len(in.ModelIDs)); err != nil {
			return nil, nil, err
		}
		rows, err := pool.Query(ctx, `
			SELECT `+modelCols+`
			FROM ev_models WHERE id = ANY($1::uuid[])`, in.ModelIDs)
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
