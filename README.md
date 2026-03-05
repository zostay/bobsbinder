# Bob's Binder

A web application that helps people collect and organize end-of-life documents for their loved ones. Be like Bob. Bob cares about his family.

## Architecture

- **Go API** — REST API built with Chi router, MySQL, JWT auth, Zap logging
- **Vue3 Frontend** — SPA built with Vuetify 3, Pinia, Vue Router, TypeScript
- **Hugo Marketing** — Landing site with Bootstrap 5, custom layouts
- **Nginx** — Reverse proxy routing all services through port 80

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose

## Getting Started

```bash
cp .env.example .env
make up
```

Then visit:

- `http://localhost/` — Marketing site
- `http://localhost/my/` — Application
- `http://localhost/api/health` — API health check

## Development

All services support hot reloading:

- **Go API** — [Air](https://github.com/air-verse/air) watches `.go` files and rebuilds automatically
- **Vue Frontend** — Vite HMR updates the browser on `.vue`/`.ts` changes
- **Hugo** — Built-in file watcher live reloads on content/layout changes

## Make Targets

```
make up              Start all services
make down            Stop all services
make restart         Restart all services
make logs            Follow logs for all services
make logs-api        Follow API logs
make logs-frontend   Follow frontend logs
make logs-marketing  Follow marketing site logs
make migrate-up      Run database migrations
make migrate-down    Rollback last migration
make migrate-new     Create a new migration file (NAME=my-migration)
make test            Run Go tests
make lint            Run Go linter
make clean           Remove containers, volumes, and temp files
make help            Show available targets
```

## Project Structure

```
bobsbinder/
├── cmd/server/          # Go API entry point
├── internal/            # Go packages (config, db, handlers, middleware, models, router)
├── migrations/          # SQL migration files
├── frontend/            # Vue3 + Vuetify SPA
├── marketing/           # Hugo marketing site
├── nginx/               # Reverse proxy config
├── docker-compose.yaml  # Service orchestration
├── Dockerfile.api       # Go API dev container
└── Makefile             # Development operations
```

## Routing

| Path | Service | Description |
|------|---------|-------------|
| `/` | Hugo (port 1313) | Marketing/landing pages |
| `/my/` | Vite (port 5173) | Vue3 application |
| `/api/` | Go (port 8080) | REST API |
