# User Management API

User management REST API in Go with MongoDB and JWT (HS256), structured as ports and adapters.

Part 2 of the challenge (Lottery Search System design) is in [`docs/lottery-search-design.md`](docs/lottery-search-design.md).

## Run

Needs Go 1.27+, Docker (Compose v2) and `make`.

```bash
# create .env with a random JWT secret
sed "s/^JWT_SECRET=.*/JWT_SECRET=$(openssl rand -hex 32)/" .env.example > .env

make up                               # API + MongoDB in Docker
curl http://localhost:8080/healthz    # {"status":"ok"}
```

For local development, `make dev` runs MongoDB in Docker and the API on the host with hot reload. `make down` stops everything.

The API docs are in [`docs/openapi.yaml`](docs/openapi.yaml), and there is an Insomnia collection in [`docs/insomnia.json`](docs/insomnia.json).

## Configuration

| Variable | Default |
|---|---|
| `HTTP_ADDR` | `:8080` |
| `LOG_LEVEL` | `info` |
| `SHUTDOWN_TIMEOUT` | `10s` |
| `MONGO_URI` | `mongodb://localhost:27017` |
| `MONGO_DB` | `user_api` |
| `JWT_SECRET` | required, 32+ bytes |
| `JWT_ISSUER` | `user-api` |
| `JWT_TTL` | `1h` |
| `USER_COUNT_INTERVAL` | `10s` |

`JWT_SECRET` has no default. The value in `.env.example` (`change-me`) is too short on purpose, so the server refuses to start until you set a real one.

## JWT

Register, log in, then send the token as a Bearer header.

```bash
curl -X POST localhost:8080/auth/register -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com","password":"password123"}'

TOKEN=$(curl -s -X POST localhost:8080/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password123"}' | jq -r .access_token)

curl localhost:8080/users -H "Authorization: Bearer $TOKEN"
```

Tokens are HS256 with the user id in `sub`, valid for `JWT_TTL`.

## API

| Method | Path | Auth | |
|---|---|---|---|
| `POST` | `/auth/register` | – | create a user |
| `POST` | `/auth/login` | – | get a token |
| `GET` | `/users` | JWT | list users |
| `GET` | `/users/{id}` | JWT | get a user |
| `PATCH` | `/users/{id}` | JWT | update your own name/email |
| `DELETE` | `/users/{id}` | JWT | delete your own account |
| `GET` | `/healthz` | – | health check |

Register:

```http
POST /auth/register
{"name":"Alice","email":"Alice@Example.com","password":"password123"}

201 Created
{"id":"6abc0fdd392b122629e56768","name":"Alice","email":"alice@example.com","created_at":"2026-09-29T19:22:05.093Z"}
```

Login:

```http
POST /auth/login
{"email":"alice@example.com","password":"password123"}

200 OK
{"access_token":"eyJhbGciOi...","token_type":"Bearer","expires_at":"2026-09-29T20:23:21Z"}
```

Update:

```http
PATCH /users/6abc0fdd392b122629e56768
Authorization: Bearer <token>
{"name":"Alice Cooper"}

200 OK
{"id":"6abc0fdd392b122629e56768","name":"Alice Cooper","email":"alice@example.com","created_at":"2026-09-29T19:22:05.093Z"}
```

`GET /users` returns an array of the same user object, and `DELETE` returns `204 No Content`.

Errors all have the same shape:

```json
{"code":"VALIDATION_FAILED","message":"validation failed","fields":{"email":"must be a valid email address","password":"must be at least 8 characters"}}
```

| Status | Code | When |
|---|---|---|
| 400 | `VALIDATION_FAILED`, `INVALID_JSON` | bad input or body |
| 401 | `MISSING_TOKEN`, `INVALID_TOKEN` | no token or bad token |
| 401 | `INVALID_CREDENTIALS` | wrong email or password |
| 403 | `FORBIDDEN` | changing someone else's account |
| 404 | `USER_NOT_FOUND` | unknown id |
| 409 | `EMAIL_TAKEN` | email already registered |
| 500 | `INTERNAL` | anything unexpected (details only in the log) |

## Tests

```bash
make test               # unit and HTTP tests, no database needed
make test-integration   # repository tests against a real MongoDB
```

## Layout

```
cmd/http/          main (wiring, graceful shutdown) and the Echo server
internal/
  authentication/  register, login
  user/            get, list, update, delete, background user counter
shared/            config, middleware, errors, logger, validation, jwt, bcrypt, mongo
docs/              OpenAPI spec, Insomnia collection, Part 2 design
```

Each module under `internal/` has `domain/` (use cases), `port/` (interfaces) and `adapter/` (Echo, Mongo, jobs).

## Notes

- Registration is how users are created. The challenge lists "user registration" and "create a new user" separately, but without roles they are the same thing, so there is only `POST /auth/register`.
- You can only update or delete your own account. Any logged-in user can list and read users.
- Emails are stored lower-cased and are unique through a MongoDB unique index.
- Passwords are bcrypt-hashed, 8–72 bytes. Login returns the same error for an unknown email and a wrong password, and takes the same time for both.
- Tokens are stateless, so deleting a user does not revoke tokens that are already out.
- `GET /users` has no pagination, since the challenge asks for "list all users".
- On SIGTERM the server drains requests, stops the counter, then closes MongoDB.
