# JobRadar

A microservices-based job search platform. It pulls remote job listings from a public API, scores each one against a configurable keyword profile, stores the results, and lets you browse the listings and track your applications on a Kanban board.

> **Status:** work in progress.

## Features

- Periodic ingestion of live remote job listings.
- Asynchronous scoring pipeline: each listing is matched against your keywords and the result is persisted.
- Job feed with a publication window filter (default 30 days, or "Anytime"), text search across role, company and description with match highlighting, score filter, and a side panel with the full (sanitized) description.
- Application tracker: create applications, list them, and move them between statuses (`applied`, `interview`, `offer`, `rejected`) on a Kanban board, using buttons or drag and drop.
- Graceful degradation: the cache and the scoring service are optional dependencies of the tracker, never blockers.

## Architecture

```mermaid
flowchart LR
    RO["RemoteOK API"] --> ING["ingestion-service"]
    ING -->|"jobs.new"| RP[("Redpanda")]
    RP -->|"jobs.new"| MAT["matching-service"]
    MAT -->|"jobs.matched"| RP
    RP -->|"jobs.matched"| JOBS["jobs-service"]
    JOBS --> PGJ[("PostgreSQL<br/>jobs_service")]
    TRK["tracker-service"] -->|"gRPC MatchScore"| MAT
    TRK --> PGT[("PostgreSQL<br/>tracker_service")]
    TRK --> RED[("Redis cache")]
    FE["Next.js frontend"] -->|"REST"| JOBS
    FE -->|"REST"| TRK
```

Three communication styles are used on purpose: REST for client-facing APIs, gRPC for direct service-to-service calls, and event streaming (Redpanda, Kafka API-compatible) for the decoupled ingestion pipeline.

### Services

| Service | Responsibility | Interfaces | Storage |
|---|---|---|---|
| `ingestion-service` | Fetches listings from RemoteOK on a timer, validates and normalizes them, publishes them. | Kafka producer: `jobs.new` | none |
| `matching-service` | Scores each listing against the configured keywords. | Kafka consumer `jobs.new`, Kafka producer `jobs.matched`, gRPC `MatchScore` | none |
| `jobs-service` | Persists scored listings and serves them. | Kafka consumer `jobs.matched`, REST | PostgreSQL (`jobs_service`) |
| `tracker-service` | Manages your applications and their status. | REST, gRPC client to `matching-service` | PostgreSQL (`tracker_service`), Redis |

Each service owns its data and its domain model. Only contracts are shared: `proto/` for gRPC and `events/` for the Kafka event payloads.

## Tech stack

| Area | Technologies |
|---|---|
| Backend | Go, standard library `net/http` |
| Messaging | Redpanda (Kafka API-compatible), `segmentio/kafka-go` |
| RPC | gRPC, Protocol Buffers |
| Data | PostgreSQL (`pgx`), Redis (`go-redis`) |
| Frontend | Next.js (App Router), TypeScript, Tailwind CSS, dnd-kit, DOMPurify, pnpm |
| Infrastructure | Docker, Docker Compose |
| Testing | `testify`, `net/http/httptest`, hand-written fakes |

## Repository layout

```
job-radar/
├── docker-compose.yml
├── go.work
├── events/                 # shared event contracts (Go module)
├── proto/                  # protobuf definitions and generated code (Go module)
├── infra/postgres/         # database and user bootstrap script
├── services/
│   ├── ingestion-service/
│   ├── matching-service/
│   ├── jobs-service/
│   └── tracker-service/
└── frontend/               # Next.js app
```

Each Go service is its own module and follows the same layering: `domain` (rules and interfaces), `repository` or `source`/`publisher` (infrastructure), `handler` or `consumer` (transport), and `cmd/main.go` as the only place where everything is wired together.

## Getting started

### Prerequisites

- Docker and Docker Compose
- Node.js 20+ and pnpm (for the frontend)
- Go and `protoc` with the Go plugins, only if you want to run tests or regenerate protobuf code

### Configuration

Configuration is read from environment variables. Each service uses its own prefix (`TRACKER_SERVICE_*`, `MATCHING_SERVICE_*`, `INGESTION_SERVICE_*`, `JOBS_SERVICE_*`).

Create two files in the repository root:

`.env` (used by Docker Compose to configure the PostgreSQL container):

```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=changeMe
```

`.env.docker` (injected into every service container). Excerpt:

```
TRACKER_SERVICE_URI=tracker:8080
TRACKER_SERVICE_DATABASE_URI=postgres://tracker_service:changeMe@postgres:5432/tracker_service?sslmode=disable
TRACKER_SERVICE_REDIS_HOST=redis_tracker
TRACKER_SERVICE_REDIS_PORT=6379
TRACKER_SERVICE_REDIS_TTL=30
TRACKER_SERVICE_MATCHING_HOST=matching:50051

MATCHING_SERVICE_GRPC_PORT=:50051
MATCHING_SERVICE_KAFKA_BROKER=redpanda:9092
MATCHING_SERVICE_PROFILE_KEYWORDS=Go,Kafka,gRPC,PostgreSQL,Redis,Docker

INGESTION_SERVICE_KAFKA_BROKER=redpanda:9092
INGESTION_SERVICE_INTERVAL_SECONDS=300

JOBS_SERVICE_URI=jobs:8081
JOBS_SERVICE_DATABASE_URI=postgres://jobs_service:changeMe@postgres:5432/jobs_service?sslmode=disable
JOBS_SERVICE_KAFKA_BROKER=redpanda:9092

FRONTEND_URI=http://localhost:3000
```

Notes:

- Do not wrap values in quotes in `.env.docker`.
- `FRONTEND_URI` is the allowed CORS origin and must match the frontend origin exactly.
- The database passwords must match the ones set in the bootstrap script in `infra/postgres/`.
- `MATCHING_SERVICE_PROFILE_KEYWORDS` is your profile: a comma-separated list of keywords.

### Run the backend

```bash
docker compose up --build
```

This starts PostgreSQL, Redis, Redpanda and the four services. The database bootstrap and the SQL migrations in `services/*/migrations` run automatically the first time PostgreSQL starts with an empty data directory.

The first ingestion cycle runs after `INGESTION_SERVICE_INTERVAL_SECONDS` seconds. For local development, set it to something small such as `30`.

### Run the frontend

```bash
cd frontend
cp .env.example .env.local
pnpm install
pnpm dev
```

`.env.local` needs:

```
NEXT_PUBLIC_JOBS_SERVICE_URL=http://localhost:8083
NEXT_PUBLIC_TRACKER_SERVICE_URL=http://localhost:8080
```

Open http://localhost:3000.

### Ports

| Component | Host port |
|---|---|
| Frontend (`pnpm dev`) | 3000 |
| `tracker-service` (REST) | 8080 |
| `jobs-service` (REST) | 8083 |
| `matching-service` (gRPC) | 50051 |
| PostgreSQL | 5432 |
| Redis | 6379 |
| Redpanda (Kafka API, external listener) | 19092 |

### Quick check

```bash
curl localhost:8083/health
curl "localhost:8083/jobs?days=7"

curl -X POST localhost:8080/applications \
  -H "Content-Type: application/json" \
  -d '{"company":"Acme","role":"Backend Engineer","url":"https://example.com/jobs/1"}'

curl -X PATCH localhost:8080/applications/1 \
  -H "Content-Type: application/json" \
  -d '{"status":"interview"}'
```

Inspect the events flowing through the pipeline:

```bash
docker compose exec redpanda rpk topic consume jobs.new -o -3
docker compose exec redpanda rpk topic consume jobs.matched -o -3
```

### Resetting everything

PostgreSQL init scripts only run when the data directory is empty. To start from scratch:

```bash
docker compose down -v
sudo rm -rf ./postgres-data
docker compose up --build
```

### Troubleshooting

- **Redpanda exits at startup on Linux or WSL2** with a message about AIO requirements: raise the kernel limit with `sudo sysctl -w fs.aio-max-nr=1048576` (persist it in your system configuration).
- **Tables or databases are missing:** the bootstrap and migrations only run on an empty data directory. See "Resetting everything".
- **`password authentication failed`:** the passwords in the bootstrap script and in `.env.docker` do not match.
- **CORS errors in the browser console:** `FRONTEND_URI` does not exactly match the frontend origin.
- **Port conflicts on 8081 or 8082:** Redpanda uses them internally, which is why `jobs-service` is published on host port 8083.
- **The feed is empty:** wait for the first ingestion cycle.

## API overview

### tracker-service (`:8080`)

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/applications` | List applications |
| POST | `/applications` | Create an application. Body: `company`, `role`, `url`, optional `matchScore` |
| PATCH | `/applications/{id}` | Change status. Body: `status` (`applied`, `interview`, `offer`, `rejected`) |

If `matchScore` is omitted, the tracker asks `matching-service` over gRPC. If that call fails, the application is still created without a score.

Response fields: `Id`, `Company`, `Role`, `OfferUrl`, `Status`, `MatchScore`, `CreatedAt`, `UpdatedAt`.

### jobs-service (`:8083` on the host)

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/jobs` | List scored jobs. Optional `days=N` returns only jobs published in the last N days. A missing or invalid value returns everything |

Response fields: `Id`, `Source`, `SourceId`, `Company`, `Role`, `OfferUrl`, `Location`, `SalaryMin`, `SalaryMax`, `Description`, `Score`, `PublishedAt`, `CreatedAt`.

### matching-service (gRPC, `:50051`)

`MatchingService.MatchScore(MatchScoreRequest) returns (MatchScoreResponse)`. The contract lives in `proto/matching/v1/matching.proto`.

To regenerate the Go code after editing the contract:

```bash
protoc \
  --go_out=proto --go_opt=module=jobradar/proto \
  --go-grpc_out=proto --go-grpc_opt=module=jobradar/proto \
  proto/matching/v1/matching.proto
```

## How scoring works

The score is the number of configured keywords that appear (case-insensitive) in the job description. It is intentionally simple.

Known limitations:

- Matching is by substring, so very short keywords (for example `Go`) can match inside unrelated words.
- The description is HTML and is matched as received.
- A raw count does not express how well a listing fits, only how many keywords it mentions.

A semantic approach (for example an LLM-based matcher behind an interface, configurable with the user's own provider key) is a planned improvement.

## Testing

Each Go service is a separate module, so run its tests from its directory:

```bash
cd services/tracker-service
go test ./... -v
```

Unit tests cover the domain rules, HTTP handlers, Kafka consumers and publishers, the RemoteOK client (against a local `httptest` server) and the repository decorators. Infrastructure is replaced with hand-written fakes behind small interfaces.

Not covered yet: the PostgreSQL repositories (they need integration tests against a real database) and the frontend.

## Design notes

- **Decorators for optional dependencies.** In `tracker-service`, the PostgreSQL repository is wrapped by a scoring decorator (gRPC) and a cache decorator (Redis, cache-aside with invalidation on writes). A failure in either is logged and tolerated, with short timeouts, so the core operation never depends on them.
- **Interfaces at the domain boundary.** Handlers and consumers depend on interfaces such as `ApplicationStore` and `JobStore`, which makes every layer testable in isolation.
- **Tolerant ingestion.** One malformed listing is skipped and logged; it never aborts a whole ingestion cycle.
- **No shared domain code.** Services only share contracts (`proto/`, `events/`). Domain models are duplicated by design to keep services independent.
- **Images built from the repository root.** The services depend on the local `proto/` and `events/` modules through `replace` directives, so the Docker build context is the repository root.

## Data source and attribution

Job listings are provided by [RemoteOK](https://remoteok.com). Their API terms require linking back to RemoteOK as the source if you display this data publicly.

## AI assistance

Developed with guidance, examples and code review from Claude Sonnet 5 (free plan, Medium effort).
