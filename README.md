# JWT Login API

A production-ready Go HTTP authentication service using **JWT (access + refresh tokens)** with Redis-backed session management, Argon2id password hashing, and PostgreSQL persistence.

## Features

- 🔐 User registration and login with JWT-based authentication
- 🔄 Access token (15 min) + refresh token (7 day) rotation
- 🍪 HttpOnly, Secure, SameSite cookie-based token delivery (XSS-safe)
- 🗂️ Redis-backed JTI cache for session validation and revocation
- 🔒 Argon2id password hashing with per-user salt (NIST recommended)
- 🐘 PostgreSQL with sqlc-generated type-safe queries
- 📦 Database migrations via `golang-migrate`

## Tech Stack

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.26 |
| **Router** | chi v5 |
| **Database** | PostgreSQL (pgx v5) |
| **Cache** | Redis (go-redis v9) |
| **Auth** | golang-jwt v5 |
| **Migrations** | golang-migrate v4 |
| **Code Generation** | sqlc |

## Prerequisites

Before running this project, ensure you have the following installed:

- [Go 1.26+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/) (for PostgreSQL and Redis)
- [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html) (for code generation)

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/improver2108/jwt-login.git
cd jwt-login
```

### 2. Configure environment variables

Create a `.env` file from the example template:

```bash
cp .env.example .env
```

Then edit `.env` with your configuration. The required variables are:

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://postgres:mysecretpassword@localhost:5432/jwtlogin?sslmode=disable` |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `ACCESS_SECRET` | Secret key for signing access tokens (base64-encoded, 32 bytes) | `QhwgEEM5qTxGhOAgnhmu4Vq7GcC8bXEBdg5jq/j+YfQ=` |
| `REFRESH_SECRET` | Secret key for signing refresh tokens (base64-encoded, 32 bytes) | `jqw3DxMJSiEBP+VX1YoRe8x02p8SeC8sgCUyMha6dIY=` |
| `APP_ENV` | Application environment (`development`, `staging`, `production`) | `development` |

> **Generating Secrets:** Generate secure keys for your environment:
> ```bash
> # Access token secret (32 bytes, base64-encoded)
> openssl rand -base64 32
> 
> # Refresh token secret (32 bytes, base64-encoded)
> openssl rand -base64 32
> ```
> Copy the output values into your `.env` file for `ACCESS_SECRET` and `REFRESH_SECRET`. Never commit secrets to version control — `.env` is listed in `.gitignore`.

### 3. Set up PostgreSQL (Docker)

Start a PostgreSQL container:

```bash
docker run -d \
  --name postgresql-jwt \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -e POSTGRES_DB=jwtlogin \
  -p 5432:5432 \
  postgres:16-alpine
```

Verify it's running:

```bash
docker exec -it postgresql-jwt psql -U postgres -d jwtlogin -c "SELECT 1;"
```

### 4. Set up Redis (Docker)

Start a Redis container:

```bash
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:latest
```

Check if it's running:

```bash
docker ps | grep redis
```

If the container was previously created but is stopped, start it with:

```bash
docker start redis
```

Test the connection via the Redis CLI:

```bash
docker exec -it redis redis-cli ping
# Expected output: PONG
```

### 5. Install Go dependencies

```bash
go mod tidy
```

## Running the Application

### Step 1: Run database migrations

Apply all pending SQL migrations to create the `users` table:

```bash
go run cmd/migration/main.go
```

This reads migration files from `db/migrations/` and applies them in order.

### Step 2: Start the API server

```bash
go run cmd/api/main.go
```

The server will start on **`http://localhost:8082`**.

## API Endpoints

| Method | Path | Description | Auth Required |
|--------|------|-------------|---------------|
| `POST` | `/register` | Register a new user | ❌ No |
| `POST` | `/login` | Log in with credentials | ❌ No |
| `GET` | `/logout` | Log out (clear session) | ✅ Yes |
| `GET` | `/refresh-token` | Refresh access token | ✅ Yes (refresh) |

### Register a new user

Creates a new user account and returns an access + refresh token pair.

**Request:**

```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword123",
  "confirm_password": "securepassword123",
  "phone": "+1234567890",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (201 Created):**

```json
{
  "message": "User registered successfully",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "phone": "+1234567890",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

Tokens are set as **HttpOnly cookies** (`access_token` and `refresh_token`) in the response.

### Log in

Authenticates a user and returns a new access + refresh token pair.

**Request:**

```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**

```json
{
  "message": "Login successful",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe"
  }
}
```

### Log out

Invalidates the current session by revoking the refresh token from Redis and clearing cookies.

**Request:**

```http
GET /logout
Cookie: access_token=...; refresh_token=...
```

**Response (200 OK):**

```json
{
  "message": "Logged out successfully"
}
```

### Refresh access token

Issues a new access token using the existing refresh token. The old session is revoked and a new one is created (token rotation).

**Request:**

```http
GET /refresh-token
Cookie: refresh_token=...
```

**Response (200 OK):**

```json
{
  "message": "Token refreshed successfully"
}
```

## Authentication Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant S as Server
    participant R as Redis
    participant D as PostgreSQL

    Note over C,D: Registration
    C->>S: POST /register (email, password)
    S->>D: Hash password (Argon2id), insert user
    S->>R: Store JTI in Redis
    S-->>C: Set access_token + refresh_token cookies

    Note over C,D: Login
    C->>S: POST /login (email, password)
    S->>D: Verify credentials
    S->>R: Store JTI in Redis
    S-->>C: Set new access_token + refresh_token cookies

    Note over C,D: Token Refresh
    C->>S: GET /refresh-token
    S->>R: Validate & revoke old JTI
    S->>R: Store new JTI
    S-->>C: Set new access_token cookie
```

## Project Structure

```
jwt-login/
├── cmd/
│   ├── api/                  # Application entrypoint (main)
│   └── migration/            # Migration runner (main)
├── db/
│   ├── migrations/           # SQL migration files (golang-migrate)
│   │   ├── 000001_create_users.up.sql
│   │   └── 000001_create_users.down.sql
│   ├── queries/              # sqlc query definitions (user.sql)
│   ├── schema/               # sqlc table schemas
│   ├── sqlc/                 # Auto-generated Go code (do not edit)
│   ├── migrate.go            # Migration logic wrapper
│   └── redis.go              # Redis client wrapper
├── internal/
│   ├── auth/                 # JWT token generation & cookie helpers
│   ├── cache/                # JTI cache abstraction (Redis-backed)
│   ├── constant/             # Shared constants (timeouts, etc.)
│   ├── errors/               # Custom error types / Sentry integration
│   ├── handler/              # HTTP request handlers
│   ├── middleware/           # Auth & validation middleware
│   ├── model/                # Domain models & DTOs
│   ├── pkg/                  # Shared utilities
│   ├── repository/           # Data access layer (PostgreSQL)
│   ├── routes/               # Route registration
│   └── service/              # Business logic layer
├── .env.example              # Environment variable template
├── .gitignore                # Excludes .env files
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
└── sqlc.yaml                 # sqlc code generation config
```

## Code Generation

This project uses [sqlc](https://sqlc.dev/) to generate type-safe Go code from raw SQL queries. This ensures database queries are validated at compile time.

After modifying queries in `db/queries/` or schema definitions in `db/schema/`, regenerate the Go code:

```bash
sqlc generate
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `connection refused` on Redis | Ensure Redis is running: `redis-cli ping` should return `PONG` |
| `could not connect to database` | Verify PostgreSQL is running and `DATABASE_URL` is correct |
| `migration error` | Run migrations before starting the API server |
| `token expired` errors | Access tokens expire after 15 minutes; use `/refresh-token` to get a new one |
| `cookie not set` in browser | Ensure cookies are enabled and SameSite policy matches your setup |

## License

This project is open source.
