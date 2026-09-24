package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golink/internal/db"
	"golink/internal/tools"
)

func main() {
	log.SetFlags(0)
	ctx := context.Background()
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := db.Migrate(pool); err != nil {
		log.Fatal(err)
	}
	if os.Getenv("SEED") != "0" {
		if err := db.SeedIfEmpty(ctx, pool); err != nil {
			log.Fatal(err)
		}
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "golink", Version: "0.1.0"}, nil)
	tools.Register(server, pool)

	if os.Getenv("TRANSPORT") == "http" || os.Getenv("PORT") != "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil)
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		mux.Handle("/", mcpHandler)
		log.Fatal(http.ListenAndServe(":"+port, mux))
	}

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
