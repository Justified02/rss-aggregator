# RSSAgg — RSS Feed Aggregator API

A production-style REST API built in Go that aggregates RSS feeds from across the web. Users register, submit feed URLs, and a concurrent background worker automatically fetches and stores posts from every feed on a schedule.

Built with real tools used in production Go shops — no ORMs, no shortcuts.

---

## Features

- **API Key Authentication** — every user gets a unique API key on registration
- **Feed Management** — submit any RSS feed URL; fetch all public feeds
- **Feed Following** — follow and unfollow feeds; get a personalised post feed
- **Concurrent Scraper** — background worker fetches all feeds in parallel on a configurable timer
- **Duplicate Prevention** — unique URL constraint ensures posts are never stored twice
- **Clean Architecture** — HTTP layer, DB layer, and scraper are fully decoupled

---

## Tech Stack

| Layer | Tool |
|---|---|
| Language | Go 1.22+ |
| Router | [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL |
| Migrations | [goose](https://github.com/pressly/goose) |
| Query Generation | [sqlc](https://sqlc.dev) |
| Environment | [godotenv](https://github.com/joho/godotenv) |

---

## Project Structure

```
rssagg/
├── main.go                  # Server setup, router, scraper bootstrap
├── handler_users.go         # User registration + auth
├── handler_feeds.go         # Feed creation + listing
├── handler_feed_follows.go  # Follow / unfollow feeds
├── handler_posts.go         # Fetch posts for authenticated user
├── middleware_auth.go       # API key extraction + user injection
├── scraper.go               # Background RSS fetcher (goroutines + WaitGroup)
├── json.go                  # Shared JSON response helpers
├── sql/
│   ├── schema/              # Goose migration files
│   └── queries/             # SQL queries (input for sqlc)
└── internal/
    └── database/            # sqlc-generated type-safe DB code
```

---

## Getting Started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/) (or Docker)
- [goose](https://github.com/pressly/goose) — `go install github.com/pressly/goose/v3/cmd/goose@latest`
- [sqlc](https://sqlc.dev) — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### 1. Clone the repo

```bash
git clone https://github.com/YOURUSERNAME/rssagg.git
cd rssagg
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` with your values:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/rssagg?sslmode=disable
```

### 3. Start PostgreSQL

Using Docker:

```bash
docker run --name rssagg \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=rssagg \
  -p 5432:5432 \
  -d postgres
```

### 4. Run migrations

```bash
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/rssagg" up
```

### 5. Run the server

```bash
go run .
```

The server starts on `http://localhost:8080`.

---

## API Reference

All authenticated routes require the header:
```
Authorization: ApiKey <your_api_key>
```

### Health Check

```
GET /v1/healthz
```

### Users

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/users` | No | Register a new user |
| `GET` | `/v1/users` | Yes | Get the authenticated user |

**Register a user:**
```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John"}'
```

```json
{
  "id": "uuid...",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "name": "John",
  "api_key": "your_api_key_here"
}
```

### Feeds

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/feeds` | Yes | Submit a new RSS feed |
| `GET` | `/v1/feeds` | No | List all feeds |

**Submit a feed:**
```bash
curl -X POST http://localhost:8080/v1/feeds \
  -H "Authorization: ApiKey <your_key>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Boot.dev Blog", "url": "https://blog.boot.dev/index.xml"}'
```

### Feed Follows

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/feed_follows` | Yes | Follow a feed |
| `GET` | `/v1/feed_follows` | Yes | Get your followed feeds |
| `DELETE` | `/v1/feed_follows/{feedFollowID}` | Yes | Unfollow a feed |

### Posts

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/v1/posts` | Yes | Get posts from feeds you follow |

Supports optional `limit` query parameter (default: 20):
```
GET /v1/posts?limit=5
```

---

## How the Scraper Works

When the server starts, a background goroutine launches alongside the HTTP server:

```
main()
  ├── go startScraping(...)   ← runs forever in background
  └── http.ListenAndServe()   ← handles API requests
```

Every minute (configurable), the scraper:

1. Fetches the N feeds that were least recently checked (`last_fetched_at ASC NULLS FIRST`)
2. Marks them as fetched immediately — preventing duplicate work
3. Spawns one goroutine per feed (all run in parallel)
4. Each goroutine fetches the RSS XML, parses it, and saves new posts
5. Duplicate URLs are silently skipped via Postgres unique constraint

This is a **work queue pattern** — the database acts as the queue, and `last_fetched_at` is the priority key.

---

## Database Schema

```
users
  id · created_at · updated_at · name · api_key

feeds
  id · created_at · updated_at · name · url · user_id · last_fetched_at

feed_follows
  id · created_at · updated_at · user_id · feed_id
  UNIQUE(user_id, feed_id)

posts
  id · created_at · updated_at · title · url · description · published_at · feed_id
  UNIQUE(url)
```

---

## Key Design Decisions

**Why sqlc over an ORM?**
sqlc generates type-safe Go functions directly from SQL. You write real SQL, get real types, and nothing is hidden behind magic. It's what production teams use when they want control.

**Why API keys instead of JWT for auth?**
RSS aggregators are used programmatically. API keys are simpler and appropriate for this use case — no token refresh, no expiry complexity.

**Why UUID over auto-increment IDs?**
UUIDs are globally unique. Safe across distributed systems, no sequential ID enumeration risk, and no ID collision if databases are ever merged.

**Why a DTO conversion layer?**
DB structs and API response shapes are kept separate. Renaming a DB column doesn't break the API contract — only the mapping function changes.

---

## What I Learned Building This

- How Go's HTTP server and `net/http` actually work under the hood
- Real SQL — migrations, foreign keys, join tables, unique constraints
- Goroutines and `sync.WaitGroup` for concurrent workloads
- Why you never expose DB structs directly to API clients (DTO pattern)
- The work queue pattern using Postgres as the queue
- Structuring a Go project that can actually be extended

---

## Roadmap

- [ ] Add pagination to posts endpoint
- [ ] Support Atom feeds (not just RSS)
- [ ] Add a `GET /v1/feeds/{feedID}/posts` endpoint
- [ ] Rate limiting on public endpoints
- [ ] Docker Compose setup for one-command local dev
- [ ] Deploy to Railway / Fly.io

---

## License

MIT
