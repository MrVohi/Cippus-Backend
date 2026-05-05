# AGENT.md — Backend

This file provides guidance to Claude Code when working in the `Backend/` directory.

## Commands

```bash
go run ./cmd/api          # start API server (default port 8080)
go run ./cmd/worker       # start AI moderation worker
go build ./...            # build all binaries
go test ./...             # run all tests
swag init -g cmd/api/main.go -o docs/  # regenerate Swagger docs (never edit docs/ manually)
```

## Folder structure

```
Backend/
├── cmd/
│   ├── api/              # main API binary (main.go)
│   └── worker/           # AI moderation worker binary (main.go)
├── internal/
│   ├── models/           # GORM model structs + migrate.go (auto-migrate list)
│   ├── handlers/         # Gin route handlers — thin, delegate to services
│   ├── services/         # business logic
│   ├── repository/       # DB queries via GORM (never raw SQL)
│   ├── middleware/       # auth, rate-limit, CORS, logging, recovery
│   └── ws/               # WebSocket hub + client
├── config/               # env loading, Config struct
├── docs/                 # Swagger output — generated, do not edit
└── migrations/           # GORM auto-migrate runs at startup; keep model list in internal/models/migrate.go
```

## API routes — prefix `/api/v1`

All responses are JSON. Protected routes require `Authorization: Bearer <access_token>`.

### Auth (public)
| Method | Path | Notes |
|--------|------|-------|
| POST | `/auth/register` | + captcha verification; returns access + refresh tokens |
| POST | `/auth/login` | returns access + refresh tokens |
| POST | `/auth/logout` | invalidates refresh token row in DB |
| POST | `/auth/refresh` | body `{refresh_token}` → new token pair |
| POST | `/auth/password-reset/request` | sends reset email; rate-limited |
| POST | `/auth/password-reset/confirm` | body `{token, new_password}` |
| GET | `/auth/oauth/:provider` | redirect to GitHub or Google |
| GET | `/auth/oauth/:provider/callback` | exchange code → merge account by email |

### Users
| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/users/:id` | public | profile + activity (posts, comments, likes) |
| PUT | `/users/me` | member | update username, avatar_url |
| GET | `/users` | admin | list all users |
| PUT | `/users/:id/role` | admin | promote/demote to moderator |
| PUT | `/users/:id/status` | admin | suspend / ban / activate |

### Posts
| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/posts` | public | list; query params: `category`, `q` (keyword search) |
| POST | `/posts` | member | creates with `pending_moderation` status; publishes to RabbitMQ |
| GET | `/posts/:id` | public | |
| PUT | `/posts/:id` | owner or mod | re-triggers moderation on content change |
| DELETE | `/posts/:id` | owner or mod | |
| POST | `/posts/:id/image` | member | multipart/form-data upload → MinIO/Azure; max 20 MB; PNG, JPEG, GIF only |

### Categories
| Method | Path | Auth |
|--------|------|------|
| GET | `/categories` | public |
| POST | `/categories` | admin |
| PUT | `/categories/:id` | admin |
| DELETE | `/categories/:id` | admin |

### Comments
| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/posts/:id/comments` | public | |
| POST | `/posts/:id/comments` | member | creates with `pending_moderation`; publishes to RabbitMQ |
| PUT | `/comments/:id` | owner or mod | |
| DELETE | `/comments/:id` | owner or mod | |

### Likes
| Method | Path | Auth | Notes |
|--------|------|------|-------|
| POST | `/posts/:id/like` | member | body `{liked: bool}`; upserts PostLike row |
| DELETE | `/posts/:id/like` | member | removes PostLike row |
| POST | `/comments/:id/like` | member | body `{liked: bool}` |
| DELETE | `/comments/:id/like` | member | |

### Notifications
| Method | Path | Auth |
|--------|------|------|
| GET | `/notifications` | member |
| PUT | `/notifications/:id/read` | member |
| PUT | `/notifications/read-all` | member |
| GET | `/notifications/preferences` | member |
| PUT | `/notifications/preferences` | member |

### Messages (DM)
| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/messages` | member | list of distinct conversation partners |
| GET | `/messages/:userId` | member | full message history with that user |
| POST | `/messages/:userId` | member | send a message (real-time delivery also via WS) |

### Reports
| Method | Path | Auth | Notes |
|--------|------|------|-------|
| POST | `/reports` | member | `{content_type: "post"\|"comment", content_id, reason}` |
| GET | `/reports` | mod | list pending reports |
| PUT | `/reports/:id` | mod | resolve or dismiss |

### AI
| Method | Path | Auth | Notes |
|--------|------|------|-------|
| POST | `/ai/improve` | member | `{text}` → improved text from Ollama; rate-limited |
| GET | `/search?q=` | public | semantic search: embed query via `nomic-embed-text`, cosine distance against `posts.embedding` (pgvector) |

### Admin — AI moderation dashboard
| Method | Path | Auth |
|--------|------|------|
| GET | `/admin/moderation` | mod |
| PUT | `/admin/moderation/posts/:id` | mod |
| PUT | `/admin/moderation/comments/:id` | mod |

### Web Push
| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/push/vapid-public-key` | public | returns VAPID public key for client subscription |
| POST | `/push/subscribe` | member | save push subscription (endpoint + keys) |
| DELETE | `/push/subscribe` | member | remove subscription |

### WebSocket
- `GET /ws?token=<access_token>` — upgrades connection; token validated before registration; used for real-time notifications and DMs.

## Middleware stack (applied in order)

1. CORS — allow `FRONTEND_URL` origin, credentials: true
2. `slog` request logger — method, path, status, latency
3. Recovery — panic → 500 + log
4. Rate limiter — per-IP on sensitive routes (see below)
5. `RequireAuth(roles...)` — injected per route group, not globally

## GORM models

**Never write raw SQL.** All queries go through GORM. Auto-migrate runs at startup; the canonical list of model structs lives in `internal/models/migrate.go`.

Conventions:
- Embed `gorm.Model` for PK + timestamps + soft delete (`DeletedAt`)
- FK fields: `UserID uint`, relation field: `User User` (pointer for optional)
- Enums are Go string constants; store as `VARCHAR` in the DB

### Model inventory

| Model | Key fields | Notes |
|-------|------------|-------|
| `User` | `Username`, `Email`, `PasswordHash` (bcrypt), `AvatarURL`, `Role` | Role: `user \| moderator \| admin` |
| `Post` | `UserID`, `Title`, `Content` (HTML/JSON), `ImageURL`, `Status` | Status: `pending_moderation \| approved \| flagged \| blocked` |
| `Category` | `Name`, `Description` | |
| `PostCategory` | `PostID`, `CategoryID` | Junction table; no `gorm.Model` needed |
| `Comment` | `UserID`, `PostID`, `Content`, `Status` | Same status enum as Post |
| `PostLike` | `UserID`, `PostID`, `Liked bool` | UNIQUE(UserID, PostID) |
| `CommentLike` | `UserID`, `CommentID`, `Liked bool` | UNIQUE(UserID, CommentID) |
| `OAuthProvider` | `UserID`, `Provider` (`github \| google`), `ProviderUserID` | |
| `RefreshToken` | `UserID`, `TokenHash`, `ExpiresAt` | |
| `PasswordResetToken` | `UserID`, `TokenHash`, `ExpiresAt`, `UsedAt *time.Time` | |
| `Notification` | `RecipientID`, `ActorID`, `Type`, `EntityType`, `EntityID uint`, `ReadAt *time.Time` | Types: `reply \| like \| mention \| dm \| post_flagged \| post_approved` |
| `Message` | `SenderID`, `ReceiverID`, `Content`, `ReadAt *time.Time` | |
| `Report` | `ReporterID`, `ContentType` (`post \| comment`), `ContentID uint` (not FK), `Reason`, `Status` | Status: `pending \| resolved \| dismissed` |
| `PushSubscription` | `UserID`, `Endpoint`, `P256dhKey`, `AuthKey` | One per browser session per user |

**Notification types:** `reply | like | mention | dm | post_flagged | post_approved`

## RabbitMQ queues

Connection string from env: `RABBITMQ_URL`.

| Queue | Publisher | Consumer | Payload |
|-------|-----------|----------|---------|
| `moderation.check` | API (on post/comment create or edit) | AI worker | `{id, content_type, content}` |
| `notifications.dispatch` | API (on like/comment/reply/DM events) | API notification service | `{user_id, notification_type, payload}` |

## AI moderation worker (`cmd/worker`)

Independent Go binary with its own Dockerfile. Scaled separately from the API in Docker Swarm.

Flow:
1. Connect to RabbitMQ, consume `moderation.check` (ack only on success)
2. For each message: POST to Ollama `/api/generate` with the content + system prompt asking for a moderation verdict
3. Parse response → `approved`, `flagged`, or `blocked`
4. Update `Post.Status` / `Comment.Status` via DB
5. If `flagged` or `blocked`, publish to `notifications.dispatch` to alert moderators

Env: `OLLAMA_BASE_URL=http://ollama:11434`, `OLLAMA_MODEL=llama3` (configurable).

## Ollama integration (API)

- **Writing assistant** (`POST /ai/improve`): POST to `OLLAMA_BASE_URL/api/generate` with the user's text + a system prompt asking for improvements. Return the improved text in the response JSON.
- **Semantic search** (`GET /search?q=`): See pgvector section below.
- **Embedding model**: `nomic-embed-text` (via `OLLAMA_BASE_URL/api/embeddings`). Separate from the generation model used for writing assistant and moderation.

## pgvector (semantic search)

PostgreSQL extension `pgvector` is required. Enable it in Docker with `ankane/pgvector` image (or `pgvector/pgvector`).

**Embedding flow:**
1. On post create/update (status → `approved`): call Ollama `/api/embeddings` with `nomic-embed-text` + the post title + content summary → store the returned vector in `posts.embedding` (GORM field type: `pgvector.Vector`).
2. Use `github.com/pgvector/pgvector-go` for Go type + GORM support.

**Search flow (`GET /search?q=`):**
1. Call Ollama `/api/embeddings` with the user's query string → get query vector.
2. Query DB: `SELECT * FROM posts ORDER BY embedding <=> $1 LIMIT 20` (cosine distance via pgvector `<=>` operator).
3. Return ranked posts.

**GORM model addition:**
```go
// Post — add embedding field
Embedding pgvector.Vector `gorm:"type:vector(768)"`
```

Create the pgvector extension in the first migration: `CREATE EXTENSION IF NOT EXISTS vector`.

Embedding dimension depends on `nomic-embed-text` output (768). Set `EMBEDDING_DIM=768` in env if you want it configurable.

## WebSocket hub (`internal/ws`)

Single hub goroutine, channel-based:
- `register chan *Client` / `unregister chan *Client`
- `broadcast chan Message` — routes to all connections matching a `UserID`

Auth: validate `?token=` access token before registering client. Use `gorilla/websocket`.

Used for:
1. Real-time in-app notifications (push from server → client)
2. Private messages (bidirectional)

## Auth

- Access token: JWT (HS256), 15 min TTL, claims: `{sub: user_id, role, email}`, signed with `JWT_SECRET`
- Refresh token: random opaque string (32 bytes hex), 7 days, stored in `RefreshTokens`
- Logout: delete `RefreshToken` row
- OAuth flow: redirect → callback → get email from provider → find or create `User` by email → upsert `OAuthProvider` row → issue token pair
- Account merging: if email already exists, attach the new provider to the existing `User` (do not create a second account)

## Rate limiting

Apply `golang.org/x/time/rate` or `ulule/limiter` per IP:

| Route | Limit |
|-------|-------|
| `POST /auth/login` | 5 req/min |
| `POST /auth/register` | 3 req/min |
| `POST /auth/password-reset/request` | 3 req/min |
| `POST /ai/improve` | 10 req/min per user |

## Environment variables

```bash
# Server
PORT=8080
FRONTEND_URL=http://localhost:3000
LOG_LEVEL=debug   # debug | info | warn | error

# Auth
JWT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=

# Database
# Docker container hostname: postgres | local dev (outside Docker): localhost
DATABASE_URL=postgres://cippus:yourpassword@localhost:5432/cippus-db

# RabbitMQ
RABBITMQ_URL=amqp://user:pass@rabbitmq:5672/

# Ollama
OLLAMA_BASE_URL=http://ollama:11434
OLLAMA_MODEL=llama3
OLLAMA_EMBED_MODEL=nomic-embed-text
EMBEDDING_DIM=768

# Object storage (MinIO locally, Azure Blob in prod)
STORAGE_ENDPOINT=
STORAGE_ACCESS_KEY=
STORAGE_SECRET_KEY=
STORAGE_BUCKET=cippus

# Web Push (VAPID)
VAPID_PUBLIC_KEY=
VAPID_PRIVATE_KEY=
VAPID_SUBJECT=mailto:admin@cippus.app

# Email (Resend — github.com/resend/resend-go/v2)
# Local dev: use onboarding@resend.dev as FROM, send to your own registered email (no domain needed)
# Production: verify your domain on Resend, then use noreply@yourdomain.com
RESEND_API_KEY=
EMAIL_FROM=onboarding@resend.dev

# Captcha (reCAPTCHA v2 "I'm not a robot" — verify token server-side via Google API)
RECAPTCHA_SECRET=
```

## Docker / Docker Swarm

`docker-compose.yml` at repo root must include `deploy.replicas` on every service to remain Swarm-compatible.

| Service | Replicas (local Swarm demo) | Stateful? |
|---------|----------------------------|-----------|
| `api` | 2+ | No |
| `worker` | 2+ | No |
| `postgres` | 1 | Yes |
| `rabbitmq` | 1 | Yes |
| `minio` | 1 | Yes (dev only) |
| `ollama` | 1 | Yes (dev only) |

Target: demonstrate live load balancing between `api` replicas (scale up/down via `docker service scale`).

## Logging

Use `log/slog` (Go stdlib ≥1.21). JSON handler in production, text handler in dev. Derive from `LOG_LEVEL` env. Never log secrets, tokens, or passwords.

## Testing

### AI moderation integration tests
- Test dataset in `internal/worker/testdata/`: messages covering normal content, insults, spam, illegal content, edge cases — each with an expected verdict
- Tests assert that Ollama's decision matches the expected verdict
- Must run in CI (GitHub Actions) on every push to `main`

## Go version

Minimum: **Go 1.25** (current `go.mod`). Required for `log/slog` stdlib, updated `net/http` patterns, and general toolchain compatibility. Pin the version in `go.mod`.

## Go conventions

- Structs, interfaces, and methods over ad-hoc functions
- Goroutines + channels for concurrency (WS hub, RabbitMQ consumer)
- Explicit error handling — no `_` ignoring errors
- MVC-style separation: handlers → services → repository
- `swag` annotations on every handler for Swagger generation
