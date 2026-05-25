# Cippus — Backend

REST API for Cippus, a community forum where users share and follow ongoing projects. Built in Go with Gin, GORM, and PostgreSQL.

## Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP framework | Gin |
| ORM | GORM + pgx (PostgreSQL driver) |
| Database | PostgreSQL + pgvector extension |
| Auth | JWT HS256 (stateless access) + DB-backed refresh tokens |
| OAuth | GitHub, Google (account merge by email) |
| Message broker | RabbitMQ |
| AI | Ollama (`llama3` for moderation/writing, `nomic-embed-text` for embeddings) |
| Semantic search | pgvector cosine distance |
| WebSocket | gorilla/websocket |
| Object storage | MinIO (local), Azure Blob Storage (production) |
| Email | Resend (`resend-go/v2`) |
| Logging | `log/slog` — JSON in production, text in dev |
| Containerization | Docker + Docker Swarm |
| CI/CD | GitHub Actions → Docker Hub → Swarm SSH deploy |

## Project structure

```
Backend/
├── cmd/
│   ├── api/              # Main API binary (main.go + routes.go)
│   └── worker/           # AI moderation worker (independent binary)
├── config/               # Env loading, Config struct, DB connection, logger setup
├── docs/                 # MCD, MLD, Swagger output (generated — do not edit)
├── internal/
│   ├── handlers/         # Gin route handlers — thin, delegate to services
│   ├── middleware/        # Auth, role guard, CORS, rate limiting, recovery
│   ├── models/           # GORM model structs + migrate.go
│   ├── repository/       # DB queries via GORM (no raw SQL)
│   ├── services/         # Business logic
│   └── ws/               # WebSocket hub + client
└── .github/workflows/    # CI/CD pipeline
```

## Getting started

### Prerequisites

- Go 1.25+
- Docker (with Swarm mode for local multi-replica testing)
- A running PostgreSQL instance with the `pgvector` extension enabled
- A running RabbitMQ instance
- [Ollama](https://ollama.ai) with `llama3` and `nomic-embed-text` pulled

### Local setup

**1. Copy the env file and fill in values:**

```bash
cp .env.example .env
```

See the [Environment variables](#environment-variables) section below for a full reference.

**2. Start infrastructure (PostgreSQL, RabbitMQ, MinIO, Ollama) via Docker Compose at the repo root:**

```bash
docker compose up -d
```

**3. Run the API server:**

```bash
go run ./cmd/api
# → listening on :8080
```

**4. (Optional) Run the AI moderation worker:**

```bash
go run ./cmd/worker
```

The worker is a separate binary that consumes the `moderation.check` RabbitMQ queue and calls Ollama to evaluate post/comment content.

### Common commands

```bash
go run ./cmd/api          # Start API server (default port 8080)
go run ./cmd/worker       # Start AI moderation worker
go build ./...            # Build all binaries
go test ./...             # Run all tests
go vet ./...              # Static analysis

# Regenerate Swagger docs (never edit docs/ manually)
swag init -g cmd/api/main.go -o docs/
```

## API overview

Base path: `/api/v1`. All responses are JSON.  
Protected routes require `Authorization: Bearer <access_token>`.

### Auth (public)

| Method | Path | Notes |
|--------|------|-------|
| POST | `/auth/register` | reCAPTCHA v2 verified server-side; returns token pair |
| POST | `/auth/login` | Returns access + refresh tokens |
| POST | `/auth/logout` | Invalidates refresh token row |
| POST | `/auth/refresh` | `{refresh_token}` → new token pair |
| POST | `/auth/password-reset/request` | Sends reset email via Resend; rate-limited |
| POST | `/auth/password-reset/confirm` | `{token, new_password}` |
| GET | `/auth/oauth/:provider` | Redirect to GitHub or Google |
| GET | `/auth/oauth/:provider/callback` | Exchange code; merge account by email |

### Users

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/users/:id` | public | Profile + activity |
| PUT | `/users/me` | member | Update username, avatar_url |
| GET | `/users` | admin | List all users |
| PUT | `/users/:id/role` | admin | Promote/demote to moderator |
| PUT | `/users/:id/status` | admin | Suspend / ban / activate |

### Posts

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/posts` | public | `?category=` `?q=` keyword filter |
| POST | `/posts` | member | Created with `pending_moderation`; enqueued to RabbitMQ |
| GET | `/posts/:id` | public | |
| PUT | `/posts/:id` | owner or mod | Re-triggers moderation |
| DELETE | `/posts/:id` | owner or mod | |
| POST | `/posts/:id/image` | member | `multipart/form-data`; max 20 MB; PNG/JPEG/GIF |

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
| POST | `/posts/:id/comments` | member | Created with `pending_moderation`; enqueued to RabbitMQ |
| PUT | `/comments/:id` | owner or mod | |
| DELETE | `/comments/:id` | owner or mod | |

### Likes

| Method | Path | Auth |
|--------|------|------|
| POST | `/posts/:id/like` | member |
| DELETE | `/posts/:id/like` | member |
| POST | `/comments/:id/like` | member |
| DELETE | `/comments/:id/like` | member |

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
| GET | `/messages` | member | List of distinct conversation partners |
| GET | `/messages/:userId` | member | Full message history |
| POST | `/messages/:userId` | member | Send; also delivered in real-time via WebSocket |

### Reports

| Method | Path | Auth |
|--------|------|------|
| POST | `/reports` | member |
| GET | `/reports` | mod |
| PUT | `/reports/:id` | mod |

### AI

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| POST | `/ai/improve` | member | `{text}` → improved text from Ollama; rate-limited |
| GET | `/search?q=` | public | Semantic search via pgvector cosine distance |

### Admin — AI moderation dashboard

| Method | Path | Auth |
|--------|------|------|
| GET | `/admin/moderation` | mod |
| PUT | `/admin/moderation/posts/:id` | mod |
| PUT | `/admin/moderation/comments/:id` | mod |

### Web Push

| Method | Path | Auth |
|--------|------|------|
| GET | `/push/vapid-public-key` | public |
| POST | `/push/subscribe` | member |
| DELETE | `/push/subscribe` | member |

### WebSocket

```
GET /ws?token=<access_token>
```

Upgrades to WebSocket after validating the access token. Used for real-time notifications and private messages.

## Authentication

- **Access token:** JWT HS256, 15-minute TTL, claims: `{sub, role, email}`, signed with `JWT_SECRET`.
- **Refresh token:** 32-byte random hex, 7-day TTL, stored hashed in `RefreshTokens`.
- **Logout:** deletes the `RefreshToken` row.
- **OAuth:** redirect → callback → look up user by email → upsert `OAuthProvider` row → issue token pair. If the email already exists, the provider is attached to the existing account — no duplicate accounts.

## Rate limits

| Route | Limit |
|-------|-------|
| `POST /auth/login` | 5 req/min per IP |
| `POST /auth/register` | 3 req/min per IP |
| `POST /auth/password-reset/request` | 3 req/min per IP |
| `POST /ai/improve` | 10 req/min per user |

## RabbitMQ queues

| Queue | Publisher | Consumer | Payload |
|-------|-----------|----------|---------|
| `moderation.check` | API (post/comment create or edit) | AI worker | `{id, content_type, content}` |
| `notifications.dispatch` | API (like/comment/reply/DM events) | API notification service | `{user_id, notification_type, payload}` |

## AI moderation worker

The worker (`cmd/worker`) is a standalone Go binary with its own Dockerfile. It scales independently from the API in Docker Swarm.

Flow:
1. Connect to RabbitMQ, consume `moderation.check` (ack only on success).
2. POST to `OLLAMA_BASE_URL/api/generate` with the content + moderation system prompt.
3. Parse verdict: `approved`, `flagged`, or `blocked`.
4. Update `Post.Status` / `Comment.Status` in the DB.
5. If `flagged` or `blocked`, publish to `notifications.dispatch` to alert moderators.

## Semantic search

Posts are embedded on creation/approval using `nomic-embed-text` via Ollama. Embeddings (768 dimensions) are stored in `posts.embedding` as a `vector(768)` column powered by the `pgvector` PostgreSQL extension.

`GET /search?q=<query>` embeds the query string and returns posts ranked by cosine distance (`<=>`).

## Database models

Auto-migrate runs at startup. All 14 models are listed in `internal/models/migrate.go`.

| Model | Key fields | Notes |
|-------|------------|-------|
| `User` | `Username`, `Email`, `PasswordHash`, `AvatarURL`, `Role` | Role: `user \| moderator \| admin` |
| `Post` | `UserID`, `Title`, `Content`, `ImageURL`, `Status`, `Embedding` | Status: `pending_moderation \| approved \| flagged \| blocked` |
| `Category` | `Name`, `Description` | |
| `PostCategory` | `PostID`, `CategoryID` | Junction table |
| `Comment` | `UserID`, `PostID`, `Content`, `Status` | Same status enum as Post |
| `PostLike` | `UserID`, `PostID`, `Liked bool` | `UNIQUE(UserID, PostID)` |
| `CommentLike` | `UserID`, `CommentID`, `Liked bool` | `UNIQUE(UserID, CommentID)` |
| `OAuthProvider` | `UserID`, `Provider`, `ProviderUserID` | `github \| google` |
| `RefreshToken` | `UserID`, `TokenHash`, `ExpiresAt` | |
| `PasswordResetToken` | `UserID`, `TokenHash`, `ExpiresAt`, `UsedAt *time.Time` | |
| `Notification` | `RecipientID`, `ActorID`, `Type`, `EntityType`, `EntityID`, `ReadAt` | Types: `reply \| like \| mention \| dm \| post_flagged \| post_approved` |
| `Message` | `SenderID`, `ReceiverID`, `Content`, `ReadAt *time.Time` | |
| `Report` | `ReporterID`, `ContentType`, `ContentID`, `Reason`, `Status` | Status: `pending \| resolved \| dismissed` |
| `PushSubscription` | `UserID`, `Endpoint`, `P256dhKey`, `AuthKey` | One per browser session |

## Environment variables

Copy `.env.example` to `.env` and fill in the values:

```bash
# Server
PORT=":8080"
FRONTEND_URL=http://localhost:3000
LOG_LEVEL=debug                  # debug | info | warn | error

# Auth
JWT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=

# Database (pgvector-enabled PostgreSQL)
DATABASE_URL=postgres://cippus:yourpassword@localhost:5432/cippus-db

# RabbitMQ
RABBITMQ_URL=amqp://cippus:yourpassword@localhost:5672/

# Ollama
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL=llama3
OLLAMA_EMBED_MODEL=nomic-embed-text
EMBEDDING_DIM=768

# Object storage — MinIO locally, Azure Blob in production
STORAGE_ENDPOINT=http://localhost:9000
STORAGE_ACCESS_KEY=
STORAGE_SECRET_KEY=
STORAGE_BUCKET=cippus

# Web Push (VAPID)
VAPID_PUBLIC_KEY=
VAPID_PRIVATE_KEY=
VAPID_SUBJECT=mailto:admin@cippus.app

# Email (Resend)
# Local dev: onboarding@resend.dev works without a verified domain
# Production: use your verified domain
RESEND_API_KEY=
EMAIL_FROM=onboarding@resend.dev

# Captcha (reCAPTCHA v2 "I'm not a robot")
RECAPTCHA_SECRET=
```

## Docker

The API has a multi-stage Dockerfile (`golang:1.25-alpine` → `alpine:3.20`, non-root user).

```bash
docker build -t cippus-api .
docker run -p 8080:8080 --env-file .env cippus-api
```

## CI/CD

GitHub Actions pipeline (`.github/workflows/CI-CD.yaml`):

1. **Build job** — `go test ./...` + `go vet ./...` on every push/PR to `main`.
2. **Deploy job** (only on `main`, after build passes):
   - Builds and pushes the Docker image to Docker Hub (tags: `:prod` + `:${{ github.sha }}`).
   - SSHes into the Swarm manager and runs `docker service update` to roll out the new image.

Required repository secrets: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`, `IMAGE_REPO`, `SWARM_HOST`, `SWARM_USER`, `SWARM_SSH_KEY`.

## Docker Swarm

The `docker-compose.yml` at the repo root is Swarm-compatible (`deploy.replicas` on every service).

| Service | Replicas | Stateful |
|---------|----------|----------|
| `api` | 2+ | No |
| `worker` | 2+ | No |
| `postgres` | 1 | Yes |
| `rabbitmq` | 1 | Yes |
| `minio` | 1 | Yes (dev only) |
| `ollama` | 1 | Yes (dev only) |

Scale at runtime: `docker service scale cippus_api=4`
