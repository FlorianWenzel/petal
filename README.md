# petal

A quiet mood journal. One petal a day, your month becomes a flower.

Stack: **Nuxt 4** (web) · **Go 1.23** (API, single binary) · **pgx + Postgres 17** · **Caddy** · **Docker Compose**.

## Layout

```
apps/
  web/          Nuxt 4 — marketing pages (SSR), journaling app (/app, SPA-only)
  api/          Go — auth (session cookies), entries CRUD, pgx, embedded migrations
infra/
  Caddyfile     reverse proxy: /api/* → api, rest → web
compose.yml     db, api, web, caddy
.env.example    copy to .env before first run
```

## Local dev (without Docker)

```bash
# 1. install web deps
npm install

# 2. start postgres only
docker compose up -d db

# 3. set api env (the Go binary auto-runs migrations on startup)
cd apps/api
echo "DATABASE_URL=postgres://petal:devpassword@localhost:5432/petal" > .env

# 4. terminals
npm run dev:api      # :4000  (cd apps/api && go run .)
npm run dev:web      # :3000  (Nitro dev middleware proxies /api → :4000)
```

## Production (Docker Compose)

```bash
cp .env.example .env
# edit .env: set POSTGRES_PASSWORD, SITE_URL, SITE_DOMAIN
docker compose up -d --build
```

Caddy gets free TLS automatically when `SITE_DOMAIN` resolves to your VPS public IP.

## Routes

- `/` — landing (prerendered)
- `/pricing`, `/about` — prerendered
- `/login`, `/signup` — SSR
- `/app/**` — SPA-only (auth required)
- `/api/auth/{signup,login,logout,me}` — session cookie auth
- `/api/auth/me/username` — rename (PATCH)
- `/api/entries/{year}/{month}` — load month
- `/api/entries/{year}/{month}/{day}` — upsert day (PUT)

## Migrations

SQL files in `apps/api/migrations/*.up.sql`, embedded into the binary via `//go:embed`. The API runs pending migrations on boot, tracked in `schema_migrations`. Existing `kysely_migration` rows are auto-imported on first Go boot.

## Backups

`pg_dump` from a cron, ship to off-server storage. Sample:

```bash
docker compose exec -T db pg_dump -U petal petal | gzip > backup-$(date +%F).sql.gz
```

## Why this stack

- **Self-hosted** — full control, predictable cost, no vendor lock-in.
- **Boring tech** — Postgres + Go + Nuxt. All long-term supported.
- **Single domain** via Caddy → cookies work without CORS games.
- **Tiny API** — single static binary, ~10MB image, ~10MB RSS. No Node runtime in prod.
- **pgx + plain SQL** — no ORM, no codegen daemon, no engine binary.
