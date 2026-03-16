# StudyBuddy Backend — Architecture

## Overview

The backend is split into **microservices**, each following **clean architecture** (domain → use cases → adapters). Services communicate over HTTP; the frontend and mobile app call service APIs (optionally through an API Gateway later).

## Principles

- **Clean architecture per service**: domain and use cases do not depend on HTTP or DB; adapters implement ports (repositories, HTTP handlers).
- **Database per service** (logical): each service owns its tables; shared DB for MVP is acceptable with clear schema ownership (e.g. `auth.users` vs `users.profiles` can live in one DB with prefixes or schemas).
- **Contract-first**: OpenAPI defines the API; backend and frontend/mobile align on it.
- **JWT for auth**: Auth service issues tokens; other services validate JWT (shared secret or JWKS) and do not call Auth on every request.

## Service boundaries (MVP)

| Service    | Responsibility                          | MVP scope                          |
|-----------|------------------------------------------|------------------------------------|
| **Auth**  | Register, login, issue/refresh JWT       | Email+password, JWT access/refresh |
| **Users** | Profile, interests, universities, degrees| User CRUD, interests, profile      |
| Courses   | Course catalog, user–course relation     | Week 3                             |
| Availability | User availability slots                | Week 3                             |
| Matching  | Candidates, match requests, invites      | Week 4                             |
| Reviews   | Reviews and ratings (or under Users)     | Week 5                             |
| Points    | Point transactions, totals (gamification)| Week 6                             |

For the first weeks, **Auth** and **Users** are enough. Add others as the plan progresses.

## What to start with

1. **Environment**: Docker Compose (Postgres), `.env`, Makefile — **done**.
2. **Documentation**: This doc + OpenAPI stub — **done**.
3. **First services**: **Auth** then **Users** (Auth has no dependency on other services; Users depends on “who is the user” via JWT).

Order of implementation:

1. **Auth service** — register, login, JWT (access + refresh). No DB in Week 1 is OK (stub in memory); add DB and migrations in Week 2.
2. **Users service** — profile CRUD, interests. Validates JWT (using shared `pkg/auth`), reads `UserID` from token.
3. Later: API Gateway (single entrypoint, route to Auth/Users), then Courses, Availability, Matching, etc.

## Clean architecture (per service)

```
service/
├── domain/           # Entities, value objects, domain errors
├── usecase/          # Application logic (input/output ports)
├── delivery/         # HTTP handlers (adapters)
├── repository/       # DB implementation (adapters)
└── (optional) client # Outgoing HTTP to other services
```

- **domain**: no imports from usecase/delivery/repository.
- **usecase**: depends only on domain + interfaces (ports) for repository and external services.
- **delivery**: implements HTTP; calls use cases.
- **repository**: implements repository interfaces; uses DB.

Shared code (JWT validation, logging, middleware) lives in **pkg/** so any service can use it without depending on another service.

## Data and DB

- **MVP**: One PostgreSQL instance; each service uses its own tables (or schema). Migrations can live in a single `/migrations` folder with naming like `001_auth_users.sql`, `002_users_profiles.sql`, or per-service migration runners later.
- **Auth** owns: users (id, email, password_hash, is_active) for login only; or Auth only issues tokens and **Users** owns the full user row — decide in Week 1 (recommendation: Auth stores only credentials; Users stores profile and links by email or internal user id).
- **Users** owns: profiles, interests, user_interests, universities, degrees.

## Security

- Passwords: bcrypt (or argon2) in Auth.
- JWT: HS256 or RS256; access token short-lived (e.g. 15 min), refresh token longer (e.g. 7 days).
- Other services: validate JWT on protected routes, extract user id/email, pass to use cases.
- CORS and rate limiting: add at gateway or per-service later.

## Next steps

1. Implement Auth: register/login endpoints, JWT issuance, stub or real DB.
2. Implement Users: profile + interests endpoints, JWT validation middleware.
3. Align OpenAPI with implemented endpoints and share with frontend/mobile.
