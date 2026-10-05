# Matrix Chat

Matrix Chat is a Go-based real-time messaging application with a terminal UI client and a PostgreSQL-backed REST/WebSocket API. It supports user registration, login, personal chat, group chat, and live message delivery over WebSockets.

The project combines:

- A Go HTTP server for authentication, messaging, and group management
- A PostgreSQL database for persistent user and chat data
- A Bubble Tea terminal client for interactive chat in a TUI
- JWT-based access tokens plus refresh-token support
- Docker support for quick local development

---

## Features

- User registration and login
- JWT access tokens with refresh-token rotation
- Direct messaging between users
- Group creation and member management
- Conversation list with latest-message preview
- Message history retrieval by conversation
- Real-time message delivery through WebSockets
- Terminal interface built with Bubble Tea
- Dockerized local development environment

---

## Tech Stack

- Go 1.26
- PostgreSQL
- Gorilla WebSocket
- SQLC for typed DB queries
- Bubble Tea + Lip Gloss for the TUI
- JWT (golang-jwt)
- Docker / Docker Compose
- Goose for schema migrations

---

## Architecture Overview

This repository is organized into a few main areas:

- `cmd/server` — starts the HTTP server
- `cmd/cli` — starts the terminal client
- `http/handlers` — API endpoint handlers
- `http/middlewares` — request auth and logging middleware
- `http/realtime` — WebSocket hub and broadcast logic
- `sql/database` — generated SQLC models and query code
- `sql/queries` and `sql/schema` — SQL queries and migrations
- `tui/` — Bubble Tea application pages and navigation logic
- `config/` — DB initialization and environment loading

Typical flow:

1. The user authenticates through the API or TUI.
2. The Go server validates JWT access tokens for protected routes.
3. Messages and metadata are persisted in PostgreSQL.
4. WebSocket clients receive real-time updates via the hub.
5. The TUI reads and displays conversation state from the server.

---

## Project Structure

```text
.
├── cmd/
│   ├── cli/
│   │   └── main.go
│   └── server/
│       └── main.go
├── config/
│   └── connect.go
├── http/
│   ├── handlers/
│   ├── helpers/
│   ├── middlewares/
│   └── realtime/
├── sql/
│   ├── database/
│   ├── queries/
│   └── schema/
├── tui/
│   ├── app/
│   ├── auth/
│   ├── client/
│   ├── pages/
│   ├── nav/
│   └── styles/
├── Dockerfile
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── go.mod
├── Makefile
├── sqlc.yml
└── .env (local, not committed)
```

---

## Prerequisites

Before running the project, make sure you have:

- Go 1.26 or newer
- Docker and Docker Compose (recommended for local database and server)
- PostgreSQL (if running without Docker)
- A `.env` file in the project root

---

## Environment Variables

Create a `.env` file at the project root with values similar to the following:

```env
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_DB=postgres
POSTGRES_PASSWORD=your_secure_password
JWT_SECRET=replace_with_a_long_secret
```

Notes:

- `POSTGRES_PASSWORD` is required for the database connection.
- `JWT_SECRET` is required for access-token signing.
- When running via Docker Compose, `POSTGRES_HOST` is typically `db`.
- If you run the server outside Docker, you may set `POSTGRES_HOST=localhost`.

---

## Running with Docker (Recommended)

### Development

```bash
docker compose -f docker-compose.dev.yml up --build
```

This starts:

- the Go API on `http://localhost:8080`
- a PostgreSQL container on `localhost:5432`

The API container uses `air` for live reload during development, as configured in the Dockerfile.

To inspect logs:

```bash
docker logs -f go-api-container
```

To stop the stack:

```bash
docker compose -f docker-compose.dev.yml down
```

---

## Running Without Docker

### Start the database

If PostgreSQL is not already available, use Docker only for the database:

```bash
docker run --name matrix-db -e POSTGRES_PASSWORD=your_secure_password -p 5432:5432 -d postgres:18.6-alpine
```

Then start the API:

```bash
go run ./cmd/server
```

Start the TUI client in another terminal:

```bash
go run ./cmd/cli
```

---

## Database Setup and Migrations

This project uses SQL migrations under `sql/schema` and generated typed SQL under `sql/database`.

To apply migrations with Goose:

```bash
make goose_up
```

To roll back:

```bash
make goose_down
```

To regenerate SQLC bindings:

```bash
make sqlc
```

You can also run Goose manually:

```bash
cd sql/schema

goose postgres "postgresql://postgres:YOUR_PASSWORD@localhost:5432/postgres?sslmode=disable" up
```

---

## Typical Project Commands

The repository includes a `Makefile` with useful shortcuts:

```bash
make run
make goose_up
make goose_down
make psql
make sqlc
```

`make run` streams the API container logs, while the other commands help with Postgres access and schema management.

---

## API Overview

The HTTP server listens on port `8080`.

### Authentication

#### Register a user

```http
POST /api/register
Content-Type: application/json
```

```json
{
  "name": "alice",
  "password": "secret123"
}
```

Response includes:

- `id`
- `name`
- `refresh_token`
- `access_token`

#### Log in

```http
POST /api/login
Content-Type: application/json
```

```json
{
  "name": "alice",
  "password": "secret123"
}
```

#### Refresh token

```http
POST /api/refresh
Content-Type: application/json
```

```json
{
  "refresh_token": "RAW_REFRESH_TOKEN"
}
```

### Protected routes

All protected routes require an Authorization header:

```http
Authorization: Bearer <access_token>
```

#### Get current user

```http
GET /api/users/me
```

#### Get user by ID

```http
GET /api/users/id/{userID}
```

#### Get user ID by username

```http
GET /api/users/{userName}
```

### Messaging

#### Send a message

```http
POST /api/conversations/{receiverID}/messages
```

Body:

```json
{
  "body": "Hello from Alice"
}
```

Note: the receiver can be either a user ID or a chat group ID.

#### Get all conversations for current user

```http
GET /api/conversations
```

#### Get messages for a conversation

```http
GET /api/conversations/{receiverID}/messages
```

### Group management

#### Create a group

```http
POST /api/groups
Content-Type: application/json
```

```json
{
  "name": "Friends"
}
```

#### Add a member to a group

```http
POST /api/groups/{groupID}/members
Content-Type: application/json
```

```json
{
  "user_id": "<USER_UUID>"
}
```

### WebSocket

```http
GET /ws
```

The server upgrades the connection if the request has a valid bearer token in the Authorization header.

---

## Real-Time Messaging

The real-time system uses a hub implementation stored in `http/realtime`.

- New WebSocket clients connect to `/ws`
- The server authenticates the user via JWT
- The hub tracks active users and sends message payloads to the relevant recipients
- Direct messages and group messages can be pushed live to connected clients

This is what enables near-instant message delivery in the client UI.

---

## Terminal Client Usage

The terminal app is started with:

```bash
go run ./cmd/cli
```

It uses Bubble Tea to render a matrix inspired chat UI with login/signup screens and messaging views.

The TUI handles:

- login and account creation
- navigation through chat pages
- listing conversations
- chatting with users or groups
- switching between menu and home flows

---

## Security Notes

- Access tokens are JWTs signed with `JWT_SECRET`.
- Refresh tokens are hashed before storage in the database.
- Passwords are hashed using bcrypt.
- Protected routes require a valid `Authorization: Bearer ...` header.
- The server validates token signature and expiry.

---

## Troubleshooting

### Database connection problems

Check that:

- PostgreSQL is running
- `.env` contains the right values
- the port matches `5432`
- the database user and password are correct

### Token errors

If authentication fails:

- verify `JWT_SECRET` is set in `.env`
- make sure the token is sent as `Authorization: Bearer <token>`
- refresh the access token if it has expired

### Docker issues

If the API does not start correctly:

```bash
docker compose -f docker-compose.dev.yml logs -f go-api-container
```

If the database is not healthy, check:

```bash
docker compose -f docker-compose.dev.yml ps
```

---

## Development Notes

- The server uses `air` in Docker for rapid rebuilds during local development.
- SQLC-generated database code lives in `sql/database` and should be regenerated when SQL queries change.
- Migrations are maintained in `sql/schema` and are applied via Goose.
- The app is intentionally split between API logic and a client UI layer, making it easier to evolve the backend and terminal frontend independently.

---

## Summary

Matrix Chat is a lightweight real-time messaging system built in Go, designed around a Postgres-backed server and a terminal-based client. It covers the full flow of authentication, direct messaging, group communication, and real-time push updates, with a simple project layout that is easy to extend.

If you want to contribute or extend the project, the most natural starting points are:

- `http/handlers` for API behavior
- `http/realtime` for live message delivery
- `tui/pages` for UI flows
- `sql/queries` and `sql/schema` for database changes
