# CLAUDE.md — Ylem Repository Guide

## What is Ylem

Ylem is an open-source real-time data streaming orchestration platform. It lets users build visual pipelines that connect Apache Kafka, Amazon SQS, Google Pub/Sub, RabbitMQ, various APIs, databases, and SaaS tools. Licensed under Apache-2.0.

## Architecture

Polyglot **microservices monorepo**. Each service runs in its own Docker container on a shared Docker network (`ylem_network`). Services communicate via **HTTP REST** and **Apache Kafka**.

### Service map

| Service | Path | Language | Port | Purpose |
|---|---|---|---|---|
| **API gateway** | `api/` | Go | 7339 | OAuth2 server for external API clients; validates tokens/scopes, proxies to pipelines/statistics |
| **Users** | `backend/users/` | Go | 7333 | User/org management, auth (JWT + Redis sessions), Google OAuth, permission checks |
| **Pipelines** | `backend/pipelines/` | Go | 7336 | Pipeline/task CRUD, scheduling, Kafka producers for task runs |
| **Integrations** | `backend/integrations/` | Go | 7337 | External connectors (Slack, Jira, SQL, BigQuery, Salesforce, Twilio, etc.) |
| **Statistics** | `backend/statistics/` | Go | 7332 | Run metrics/logs stored in ClickHouse, Kafka consumer for results |
| **Task Runner** | `processor/taskrunner/` | Go | 7335 | Kafka consumer that executes tasks; load balancer partitions work |
| **Python Processor** | `processor/python_processor/` | Python | 7338 | FastAPI service running user-written Python code via RestrictedPython |
| **Nginx** | `server/` | — | 7331 | Reverse proxy, CORS, blocks private endpoints from external access |
| **UI** | `ui/` | JavaScript | 7330 (dev) | React SPA — pipeline editor, dashboards |
| **Database** | `database/` | SQL | 3306 | MySQL bootstrap (`init.sql`), Redis config |

### Infrastructure dependencies

- **MySQL** — shared instance, per-service logical databases (`users`, `api`, `integrations`, `pipelines`)
- **Redis** — session storage (users service), task runner helpers
- **Apache Kafka** — pipeline task execution message bus (Confluent image, bundled)
- **ClickHouse** — statistics service only (dedicated container)

### Kafka topics and data flow

```
pipelines (schedule/trigger)
  → task_runs
  → [load balancer] → task_runs_load_balanced
  → [task runner] → task_run_results
                   → query_task_run_results
                   → notification_task_run_results
  → [statistics resultlistener] writes to ClickHouse
```

### HTTP routing through Nginx

Nginx on port 7331 maps prefixes to services:
- `/user-api/` → `ylem_users:7333`
- `/pipeline-api/` → `ylem_pipelines:7336`
- `/integration-api/` → `ylem_integrations:7337`
- `/stats-api/` → `ylem_statistics:7332`
- `/oauth-api/` → `ylem_api:7339`

Private endpoints (`/user-api/private/`, `/integration-api/private/`) are blocked at the Nginx layer.

### Cross-service auth flow

1. Users service issues JWT tokens with Redis-backed session IDs
2. Other backend services validate tokens by calling `POST /private/user/check-authorization` on the users service
3. The API gateway uses OAuth2 with scopes for external API clients

## Repository structure

```
ylem/
├── .github/workflows/       # CI: go.yml (build+test), golangci-lint.yml
├── api/                     # Go module: ylem_api
├── backend/
│   ├── integrations/        # Go module: ylem_integrations
│   ├── pipelines/           # Go module: ylem_pipelines
│   ├── statistics/          # Go module: ylem_statistics
│   └── users/               # Go module: ylem_users
├── database/                # MySQL init.sql, data volume, Redis config
├── processor/
│   ├── python_processor/    # Python FastAPI service
│   └── taskrunner/          # Go module: task runner + load balancer
├── server/                  # Nginx config
├── ui/                      # React SPA (Create React App)
├── .env                     # DB/Redis root credentials
├── .env.common              # Shared service URLs, Kafka config, integration secrets
├── docker-compose.yml       # Root orchestration (includes per-service compose files)
├── Dockerfile-db            # MySQL image
└── Dockerfile-kafka         # Custom Kafka image (if needed)
```

Each Go service follows a similar internal layout:
```
service/
├── main.go                  # Entry: urfave/cli app with graceful shutdown
├── go.mod / go.sum
├── cli/command/             # CLI subcommands: server, migrations, Kafka workers
├── config/                  # Config via godotenv + envconfig
├── api/ or services/server/ # HTTP handlers (Gorilla mux)
├── entities/ or model/      # Domain types
├── repositories/            # Data access (GORM / raw SQL)
├── services/                # Business logic
├── db/migration/            # SQL migration files (golang-migrate)
├── tests/                   # Unit tests
├── Dockerfile
└── docker-compose.yml
```

## Language and framework details

### Go services (Go 1.27.4)

- **HTTP**: `gorilla/mux` — flat REST routes per service
- **CLI**: `urfave/cli/v2` — subcommands for `server serve`, `db migrations migrate`, Kafka workers
- **ORM**: GORM (MySQL for most services, ClickHouse driver for statistics)
- **Kafka**: `lovoo/goka` (high-level Kafka processor/emitter), IBM Sarama underneath
- **Config**: `joho/godotenv` + `kelseyhightower/envconfig` with struct tags
- **Migrations**: `golang-migrate/migrate/v4` with `file://` source, `000001_name.{up,down}.sql` naming
- **Auth**: `golang-jwt/jwt/v4`, `go-oauth2/oauth2` (API gateway)
- **Logging**: `sirupsen/logrus` (imported as `log`), level set from `YLEM_LOG_LEVEL`
- **HTTP client**: `go-resty/resty/v2` for inter-service calls

### Python processor (Python 3.12)

- FastAPI + Uvicorn
- RestrictedPython for sandboxed user code execution
- pandas for data manipulation

### Frontend (React 18)

- Create React App (`react-scripts` 5)
- Redux for state management
- MUI + Ant Design + Bootstrap (mixed component libraries)
- react-flow for pipeline graph editor
- Chart.js for visualizations
- axios for HTTP, proxied through Nginx in dev (`"proxy": "http://ylem_server:7331/"`)

## Development workflow

### Running locally

```bash
# Build and start everything from source
docker compose up

# Rebuild a specific container
docker compose build --no-cache ylem_users

# Run in background
docker compose up -d
```

The UI is available at `http://127.0.0.1:7330/` (dev) or `http://127.0.0.1:7331/` (through Nginx).

### Building individual Go services

Each Go service is an independent module. Build/test from the repo root:
```bash
go build -C api
go test -C api ./tests/...

go build -C backend/pipelines
go test -C backend/pipelines ./tests/...
```

### Linting

```bash
# Per-module, from the service directory
golangci-lint run
```

CI uses golangci-lint v2.15.0 with `--timeout=5m`.

### Running migrations

Migrations run automatically as one-shot Docker containers before services start. To run manually, each service exposes a CLI command:
```bash
# Inside the service container
./binary db migrations migrate
```

## Configuration

All configuration is environment-variable-driven:

- **`.env`** — database root credentials, Redis password
- **`.env.common`** — shared cross-service config: DB host/port, Redis host, Kafka bootstrap servers, topic names, inter-service base URLs, optional integration secrets (Twilio, AWS KMS)
- **Per-service `.env`** — service-specific database name, listen address

Config is loaded via `godotenv.Load()` then parsed into typed structs with `envconfig.Process()` using service-specific prefixes (`API`, `STATISTICS`, `TASK_RUNNER`, etc.).

Key environment variables:
- `YLEM_DATABASE_*` — MySQL connection
- `YLEM_REDIS_*` — Redis connection
- `YLEM_KAFKA_BOOTSTRAP_SERVERS` — Kafka brokers
- `YLEM_KAFKA_*_TOPIC` — Kafka topic names
- `YLEM_*_BASE_URL` — inter-service HTTP base URLs
- `YLEM_LOG_LEVEL` — logging verbosity (default: `trace`)

## CI/CD

GitHub Actions on push/PR to `main`:

1. **`go.yml`** — Builds and tests: `api`, `backend/pipelines`, `backend/users`, `backend/statistics`. Note: `backend/integrations` and `processor/taskrunner` are currently **commented out**.
2. **`golangci-lint.yml`** — Lints the same four services. Integrations and taskrunner lint also commented out.

No automated deployment pipeline is configured in the repo.

## Code conventions

### Go

- Each service is a separate Go module (no Go workspace file)
- Package names are lowercase; internal packages grouped by role (`config`, `api`, `services`, `entities`, `helpers`)
- Handler files use mixed naming: some PascalCase (`LoginUser.go`), some snake_case (`create_client.go`)
- Error handling: structured JSON error responses via helper functions; uses `errors.Is`/`errors.As` for typed errors; `fmt.Errorf` for wrapping
- Tests use stdlib `testing` only (no testify); table-driven with `t.Run`; tests live under `tests/` directories
- Editor config: tabs for indentation, final newline (`.editorconfig` per module)

### Git commits

Short imperative subject lines, sometimes with PR references:
```
Added API integration with Gemini (#3)
Fix i/o between pipelines
Kafka container and demo pipeline fixes
```

## Shared private modules

Several private Go modules under `github.com/ylem-co/` are used across services:
- `shared-messaging` — Kafka message type definitions
- `es-sql-client`, `hubspot-client`, `salesforce-client`, `opsgenie-client` — integration-specific clients

## Key directories to ignore

From `.gitignore`: `database/data/`, `the-container/`, `.idea/`, `upd_ui/`, Go workspace file (`go.work`).

`the-container/` contains consolidated Docker packaging with its own `.git` history — it's a distribution artifact, not part of the main development flow.

`upd_ui/` is a secondary/experimental React app, separate from the main `ui/`.
