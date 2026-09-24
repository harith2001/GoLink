package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golink/internal/db"
)

func TestSearchModelsLive(t *testing.T) {
	ctx := context.Background()
	pool, err := db.Connect(ctx)
	if err != nil {
		t.Skip(err)
	}
	defer pool.Close()

	server := mcp.NewServer(&mcp.Implementation{Name: "golink-test", Version: "0"}, nil)
	Register(server, pool)
	cT, sT := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, sT, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0"}, nil)
	session, err := client.Connect(ctx, cT, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "search_models",
		Arguments: map[string]any{"max_budget_lkr": 15000000, "min_range_km": 300},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("tool error: %v", res.Content)
	}
	rows, ok := res.StructuredContent.([]any)
	if !ok || len(rows) == 0 {
		t.Fatalf("expected models, got %#v", res.StructuredContent)
	}
	first := rows[0].(map[string]any)

	for _, call := range []mcp.CallToolParams{
		{Name: "get_model_details", Arguments: map[string]any{"make": first["make"], "model": first["model"]}},
		{Name: "compare_models", Arguments: map[string]any{"model_ids": []string{
			"00000000-0000-4000-8000-000000000004",
			"00000000-0000-4000-8000-000000000008",
		}}},
		{Name: "get_price_trend", Arguments: map[string]any{"model_id": "00000000-0000-4000-8000-000000000004", "months": 24}},
		{Name: "find_charging_stations", Arguments: map[string]any{"near_location": "Colombo", "fast_charging_only": true}},
		{Name: "get_import_policy", Arguments: map[string]any{}},
	} {
		out, err := session.CallTool(ctx, &call)
		if err != nil {
			t.Fatal(call.Name, err)
		}
		if out.IsError {
			t.Fatalf("%s: %v", call.Name, out.Content)
		}
		if out.StructuredContent == nil {
			t.Fatalf("%s: empty", call.Name)
		}
	}
}
