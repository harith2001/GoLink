# GoLink

An AI-assisted explorer for Sri Lanka's EV market. A Go MCP server exposes curated model specs, asking prices, charging stations, and import-duty history. Claude calls those tools. A Next.js chat UI is the planned front door (not built yet).

Ask things like:

- "EVs under 15M LKR with 300km+ range"
- "Compare the MG ZS EV and BYD Atto 3"
- "What's the latest import duty?"

## Why

- Real `pgx` / Postgres work in Go.
- A deployed example of the Anthropic Messages API MCP connector.
- A useful personal tool, not a toy CRUD demo.
- One story that covers a database, a tool-calling backend, and (once Phase 7 lands) a chat frontend.

## Status

| Piece | State |
|---|---|
| Postgres schema, seed, six read-only MCP tools | Working locally |
| stdio transport (Claude Desktop) | Working; you still add the config below |
| Streamable HTTP + Docker image | Working locally. `GET /health` returns `ok` |
| Railway (public URL + hosted Postgres) | Not deployed |
| Next.js chat UI | Not started |
| Auth, rate limit, read-only DB role | Not started |

Every tool is read-only. Nothing in v1 writes to the database.

## Architecture

```mermaid
flowchart LR
  UI["Next.js chat UI"] --> API["Next.js API route"]
  API --> Claude["Anthropic Messages API"]
  Claude --> Go["Go MCP server, pgx"]
  Go --> PG[("Postgres")]
```

Requests go down that chain. Answers come back up.

Locally, Claude Desktop talks to the server over **stdio**. In production the same server listens with **Streamable HTTP** so Claude can reach it through the API's `mcp_servers` parameter. The frontend is not in the repo yet; the Go server and database are.

## Tech stack

| Layer | Choice |
|---|---|
| MCP server | Go 1.22+ (developed on 1.26), [`github.com/modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk), [`pgx/v5`](https://github.com/jackc/pgx) |
| Database | Postgres 16. Docker locally, Railway planned |
| Migrations | [goose](https://github.com/pressly/goose) |
| Frontend | Next.js App Router (Phase 7) |
| LLM | Anthropic Messages API, MCP connector (`mcp_servers` + beta header) |
| Deploy | Railway for the Go server and Postgres (`Dockerfile`, `railway.toml`). Vercel planned for the frontend |

## MCP tools

| Tool | Input | Output |
|---|---|---|
| `search_models` | `max_budget_lkr?`, `min_range_km?`, `body_type?`, `import_type?` | Matching models: make, model, year, price range, range |
| `get_model_details` | `make`, `model` | Full spec plus min / avg / max from `price_listings` |
| `compare_models` | `model_ids` (2 or 3 UUIDs) | Side-by-side specs |
| `get_price_trend` | `model_id`, `months?` (default 12, cap 120) | Listings for that model in the window |
| `find_charging_stations` | `near_location?`, `connector_type?`, `fast_charging_only?` | Stations. `near_location` is a case-insensitive substring of `location_name` |
| `get_import_policy` | `as_of_date?` (`YYYY-MM-DD`, omit for the latest row) | Duty rate and notes in effect on that date |

`body_type` is `hatchback`, `sedan`, `suv`, or `crossover`. `import_type` is `brand_new` or `used`. `connector_type` is `CCS2`, `CHAdeMO`, or `Type2`.

Queries are parameterized (`$1`, `$2`, …). No string-built SQL.

## Database

Four tables, created by `backend/migrations/00001_init.sql`:

- **`ev_models`** — make, model, year, battery kWh, range, asking-price range, body type, import type
- **`price_listings`** — a listing for one model (source, price, date, condition, mileage). `model_id` references `ev_models` and cascades on delete
- **`charging_stations`** — name, operator, place, lat/long, connector, charger count, fast-charging flag
- **`import_policy`** — effective date, duty percent, description

Seed files live in `seed/` and load once, when `ev_models` is empty:

| File | Rows |
|---|---|
| `ev_models.sql` | 32 models |
| `price_listings.sql` | 20 listings |
| `charging_stations.sql` | 12 stations |
| `import_policy.sql` | 5 policy rows |

Prices, station coordinates, and duty rates are **curated demo figures**, not a live scrape and not a Customs gazette extract. Check a gazette before relying on a duty number. Set `SEED=0` to skip loading.

## Repo layout

```
GoLink/
  backend/
    cmd/server/main.go      # migrate, seed, stdio or HTTP
    internal/db/            # pool, goose, seed loader
    internal/models/        # row structs
    internal/tools/         # one file per MCP tool
    migrations/00001_init.sql
  seed/                     # SQL inserts
  docker-compose.yml        # local Postgres on host port 5433
  Dockerfile                # Railway / local image
  railway.toml
  .env.example
  project.md                # build plan
  GoLink-project-overview.md
```

`frontend/` arrives in Phase 7.

## Prerequisites

- Go 1.22 or newer
- Docker, for local Postgres
- Optional: [Claude Desktop](https://claude.ai/download) to call the tools over stdio
- Later: an [Anthropic API key](https://console.anthropic.com/) and a [Railway](https://railway.app/) account

## Run locally

Host port **5433** is used on purpose. A Windows Postgres install often already owns 5432.

```powershell
copy .env.example .env
docker compose up -d
cd backend
go run ./cmd/server
```

With no `PORT` and no `TRANSPORT=http`, the process speaks MCP on stdin/stdout. It applies goose migrations, then loads seed data if the models table is empty.

`DATABASE_URL` defaults to:

```
postgres://golink:golink@127.0.0.1:5433/golink?sslmode=disable
```

Use `127.0.0.1`, not `localhost`, so the client does not hit an IPv6 Postgres on 5432.

### Claude Desktop (stdio)

Add this under `mcpServers` in `%APPDATA%\Claude\claude_desktop_config.json`, then restart Claude Desktop:

```json
{
  "mcpServers": {
    "golink": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "cwd": "C:\\Users\\DELL\\Desktop\\Idea\\GoLink\\backend",
      "env": {
        "DATABASE_URL": "postgres://golink:golink@127.0.0.1:5433/golink?sslmode=disable"
      }
    }
  }
}
```

Change `cwd` if the repo lives somewhere else. A question such as "EVs under 15M LKR with 300km+ range" should come back from `search_models`.

### Streamable HTTP

```powershell
$env:TRANSPORT = "http"
$env:PORT = "8090"
$env:DATABASE_URL = "postgres://golink:golink@127.0.0.1:5433/golink?sslmode=disable"
go run ./cmd/server
```

- MCP endpoint: `http://127.0.0.1:8090/`
- Health: `http://127.0.0.1:8090/health` → `ok`

`PORT` alone also selects HTTP (Railway sets it). If `PORT` is empty and `TRANSPORT=http`, the server listens on **8080**.

The image does the same thing:

```powershell
docker build -t golink-server .
docker run --rm -p 8091:8090 `
  -e TRANSPORT=http -e PORT=8090 `
  -e DATABASE_URL=postgres://golink:golink@host.docker.internal:5433/golink?sslmode=disable `
  golink-server
```

Health is then `http://127.0.0.1:8091/health`.

### Environment

| Variable | Role |
|---|---|
| `DATABASE_URL` | Postgres URL. Default is the local Docker database on port 5433 |
| `TRANSPORT` | `http` forces Streamable HTTP. Otherwise stdio, unless `PORT` is set |
| `PORT` | HTTP listen port. Railway sets this |
| `SEED` | `0` skips the seed load. Any other value loads seed when `ev_models` is empty |
| `MIGRATIONS_DIR` | Default `migrations` (relative to the working directory) |
| `SEED_DIR` | Default `../seed` from `backend/`. The image sets `/app/seed` |
| `ANTHROPIC_API_KEY` | Frontend, Phase 7. Unused by the Go server |
| `MCP_SERVER_URL` | Public server URL, after Railway |
| `MCP_AUTH_TOKEN` | Planned for Phase 8. Not checked yet |

### Tests

From `backend/`:

```powershell
$env:DATABASE_URL = "postgres://golink:golink@127.0.0.1:5433/golink?sslmode=disable"
go test ./...
```

`TestSearchModelsLive` calls all six tools against that database. It skips if Postgres is down.

## Deploy (Railway)

Not done yet. The image is ready:

1. Create a Railway project with Postgres and copy its `DATABASE_URL`.
2. Deploy this repo with the root `Dockerfile` (`railway.toml` points at it).
3. Set `DATABASE_URL` on the Go service. Leave `SEED` unset for the first boot so the empty database gets the seed files. Railway sets `PORT`.
4. Open `https://<your-host>/health` and expect `ok`.
5. Put that origin in `MCP_SERVER_URL` for the frontend.

On startup the server runs goose, then seeds if `ev_models` has no rows.

## What is left

1. **Frontend** — Next.js chat UI. An API route calls the Messages API with `mcp_servers` aimed at the Railway URL and the MCP beta header, and streams the reply.
2. **Harden** — bearer token on the MCP endpoint, basic rate limiting, a read-only Postgres role. Queries stay parameterized.
3. **Ship** — public URL, short writeup, screenshots once the UI exists.

### Later (v2)

- Replace the curated seed with a scheduled scrape (riyasewana, patpat) into `price_listings`.
- A Docker/Kubernetes manifest for the Go server, as CKAD practice. Railway stays the real host.
- A write tool (for example "save this model to my shortlist") only after auth exists.

## Data note

This is a demo dataset for learning the MCP loop. Dealer prices move, and import duty is whatever the current gazette says. The descriptions in `import_policy` say the same thing.
