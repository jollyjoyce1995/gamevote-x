# GameVote

A small real-time app for picking which game a group plays: nominate games
(searched on Steam), vote, see the results. It also counts beers per person.
Built for a few friends, so it is designed to idle on almost nothing.

## Stack

| Part     | Tech |
|----------|------|
| Backend  | Go (stdlib `net/http`, `pgx`), one static binary |
| Database | PostgreSQL, tuned for a tiny workload (see `docker-compose.yml`) |
| Frontend | Vue 3 + TypeScript + Vite, served by the Go binary |
| API      | OpenAPI 3 spec -> generated Go server stubs **and** generated TypeScript client |
| Live updates | Server-Sent Events |

Footprint: roughly 25-40 MB for Postgres plus 15-25 MB for the app.

## Run it

```bash
cp .env.example .env        # optional: set POSTGRES_PASSWORD
docker compose up -d --build
```

Open http://localhost:8080. Data lives in the `pgdata` volume.

## Develop

```bash
docker compose up -d db                     # Postgres on 127.0.0.1:5432
cd gamevote-api-go && go run ./cmd/api      # API on :8080 (creates the schema itself)
cd gamevote-ui && bun install && bun run dev  # UI on :5173, proxies /api to :8080
```

Tests run against a real Postgres:

```bash
cd gamevote-api-go
DATABASE_URL=postgres://gamevote:gamevote@localhost:5432/gamevote?sslmode=disable go test ./...
```

## The API contract

`gamevote-api-go/api/openapi.yaml` is the single source of truth. After editing it:

```bash
cd gamevote-api-go && go generate ./api     # Go server types + handler interface
cd gamevote-ui && bun run generate          # TypeScript client in src/generated-api
```

Both generated outputs are committed. The Go compiler then tells you which handlers
to update, and `bun run type-check` tells you which frontend calls to fix.

| Operation | Endpoint |
|-----------|----------|
| `login` | `POST /api/login` - register-or-login, also used to validate a stored username |
| `listUsers` | `GET /api/users` |
| `listDrinks` / `createDrink` | `GET` / `POST /api/drinks` |
| `searchGames` | `GET /api/games?q=` - live Steam store search (cached 10 min) |
| `listParties` / `createParty` | `GET` / `POST /api/parties` |
| `getParty` | `GET /api/parties/{code}` |
| `setPartyStatus` | `PUT /api/parties/{code}/status` - also starts the next round (RESULTS -> NOMINATION) |
| `joinParty` | `PUT /api/parties/{code}/attendees/{username}` |
| `addOption` / `removeOption` | `POST /api/parties/{code}/options`, `DELETE .../options/{name}` |
| `submitVotes` | `PUT /api/parties/{code}/votes/{username}` |
| `addBeer` | `POST /api/parties/{code}/beers` |
| `streamParty` | `GET /api/parties/{code}/stream?username=` (SSE: `party_updated`, `online_users`) |

## How a party works

`NOMINATION` (add games) -> `VOTING` (-1 / 0 / +1 per game) -> `RESULTS`. The party moves
to `RESULTS` by itself when the last attendee votes. From `RESULTS`, "new round" goes back
to `NOMINATION` with the round counter increased; earlier rounds' results stay visible.

## Notes

- Identity is just a username kept in a cookie - no passwords. Anyone who can reach
  the server can act as anyone. Fine for friends on a LAN or behind a VPN; put real
  auth in front of it before exposing it to the internet.
- Config is by environment variable: `DATABASE_URL`, `PORT`, `STATIC_DIR`, `LOG_LEVEL`, `LOG_FORMAT`.
- The schema is in `gamevote-api-go/internal/store/schema.sql` and is applied (idempotently) on start.
