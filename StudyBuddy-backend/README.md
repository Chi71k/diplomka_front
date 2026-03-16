# StudyBuddy Backend

Backend for StudyBuddy — a platform for student matchmaking and collaboration (diploma project). Built with **Go**, **microservices**, and **clean architecture**.

## Stack

- **Go 1.21+**
- **PostgreSQL 16**
- **Docker** & **Docker Compose** for local dev

## Quick start

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose

### 1. Clone and prepare environment

```bash
git clone <repo-url>
cd StudyBuddy-backend
cp .env.example .env
# Edit .env if needed (JWT_SECRET, ports)
```

### 2. Start infrastructure

```bash
docker compose up -d
make db-wait   # optional: wait for Postgres ready
```

### 3. Run migrations

Install the migrate CLI (one-time):

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Then run:

```bash
make migrate-up
# Or: migrate -path migrations -database "postgres://studybuddy:studybuddy@localhost:5432/studybuddy?sslmode=disable" up
```

### 4. Run services

```bash
# Terminal 1 – Auth
make run-auth

# Terminal 2 – Users
make run-users
```

Or run from IDE: set working directory to repo root and run `cmd/auth` and `cmd/users`.

### 5. Health checks

- Auth: `GET http://localhost:8080/health`
- Users: `GET http://localhost:8081/health`

## Project layout

```
StudyBuddy-backend/
├── cmd/
│   ├── auth/               # Auth service entrypoint
│   └── users/              # Users service entrypoint
├── pkg/                    # Shared libraries (no service deps)
│   ├── auth/               # JWT issue/validate, middleware, context
│   ├── httputil/           # JSON response helpers
│   └── password/           # bcrypt hash/compare
├── services/
│   ├── auth/               # Auth: domain, usecase, delivery, repository
│   └── users/              # Users: domain, usecase, delivery, repository
├── docs/
│   ├── architecture.md
│   └── openapi/
├── migrations/           # SQL migrations (golang-migrate)
│   ├── 000001_create_users_table.up.sql
│   ├── 000002_create_interests_and_user_interests.up.sql
│   └── *.down.sql
├── docker-compose.yml
├── Makefile
└── go.mod
```

### Migrations

| Version | Description |
|--------|-------------|
| 000001 | `users` table (id, email, password_hash, first_name, last_name, bio, avatar_url, is_active, created_at, updated_at). Used by Auth and Users services. |
| 000002 | `interests` catalog and `user_interests` junction table; seeds default interests (Programming, Mathematics, etc.). |

## Development plan (MVP)

| Week | Focus |
|------|--------|
| 1 | Foundation, OpenAPI, Docker, Auth + Users stubs |
| 2 | Real auth (JWT), Users CRUD, Interests |
| 3 | Courses, Availability (manual slots) |
| 4 | Matching v1, MatchRequests, MatchInvites |
| 5 | Reviews, ratings, reputation |
| 6 | Points (transactions), minimal gamification |
| 7 | Polish, testing, deployment |
| 8 | Buffer, demo, documentation |

## Documentation

- [Architecture](docs/architecture.md) — services, boundaries, clean architecture
- [OpenAPI](docs/openapi/) — API contract (align with frontend/mobile)

## License

Diploma project — internal use.
