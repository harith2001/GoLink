package tools

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Register(s *mcp.Server, pool *pgxpool.Pool) {
	registerSearch(s, pool)
	registerDetails(s, pool)
	registerCompare(s, pool)
	registerTrend(s, pool)
	registerStations(s, pool)
	registerPolicy(s, pool)
}
