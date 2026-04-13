# StudyBuddy backend - development commands
# Use: make <target>

.PHONY: run-auth run-users run-courses build test docker-up docker-down db-wait migrate-up tidy

# Default: show help
help:
	@echo "Targets: run-auth, run-users, build, test, docker-up, docker-down, db-wait, migrate-up, tidy"

# ---------- Run services (from repo root, load .env yourself or export) ----------
run-auth:
	go run ./cmd/auth

run-users:
	go run ./cmd/users

run-courses:
	go run ./cmd/courses

# ---------- Build ----------
build:
	go build -o bin/auth ./cmd/auth
	go build -o bin/users ./cmd/users
	go build -o bin/courses ./cmd/courses

# ---------- Test ----------
test:
	go test ./...

# ---------- Docker ----------
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Wait for Postgres to be ready (requires pg_isready or use a retry loop)
db-wait:
	@echo "Waiting for Postgres..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if docker compose exec -T postgres pg_isready -U studybuddy -d studybuddy 2>/dev/null; then \
			echo "Postgres is ready"; exit 0; \
		fi; \
		sleep 2; \
	done; \
	echo "Postgres did not become ready"; exit 1

# ---------- Migrations (requires: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest) ----------
DATABASE_URL ?= postgres://user:password@localhost:5432/studybuddy?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down -all

# ---------- Deps ----------
tidy:
	go mod tidy
