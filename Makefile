.PHONY: dev mongo run test test-integration cover up down logs

# Step 1: load .env when it exists so its values reach go run and air without being printed
-include .env
export

# Step 2: hot reload, MongoDB in docker and the API rebuilt by air on every .go change
dev: mongo
	@go tool air

# Step 3: start only MongoDB in docker and wait until it is healthy
mongo:
	@docker compose up -d --wait mongo

# Step 4: run the API once without hot reload
run: mongo
	@go run ./cmd/http

# Step 5: unit and HTTP tests with the race detector, no database needed
test:
	@go test -race ./...

# Step 6: repository tests against a real MongoDB
test-integration: mongo
	@go test -race -count=1 -tags integration ./internal/user/adapter/outbound/mongo/ ./internal/authentication/adapter/outbound/mongo/

# Step 7: total coverage across every package
cover:
	@go test -coverpkg=./... -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1

# Step 8: build and start API and MongoDB fully in docker
up:
	@docker compose up --build -d --wait

# Step 9: follow the API container logs
logs:
	@docker compose logs -f api

# Step 10: stop the stack, data stays in the named volume
down:
	@docker compose down
