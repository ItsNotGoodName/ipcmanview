export DB_PATH=sqlite.db
export VITE_HOST=127.0.0.1

-include .env

migrate:
	goose -dir migrations/sql sqlite3 "$(DB_PATH)" up

clean:
	rm ${DB_PATH}
	rm ${DB_PATH}-shm
	rm ${DB_PATH}-wal

# Preview

preview:
	cd internal/web && pnpm run build && cd ../.. && go run ./cmd/ipcmanview serve

# Run

run:
	go run ./cmd/ipcmanview serve

# Dev

dev:
	air

dev-assets:
	cd internal/web && pnpm install && pnpm run dev

# Gen

gen: gen-sqlc

gen-sqlc:
	sqlc generate

# Database

db-inspect:
	atlas schema inspect --env local

db-migration:
	atlas migrate diff $(name) --env local

# Tooling

tooling: tooling-air tooling-task tooling-goose tooling-atlas tooling-sqlc

tooling-air:
	go install github.com/cosmtrek/air@latest

tooling-task:
	go install github.com/go-task/task/v3/cmd/task@latest

tooling-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

tooling-atlas:
	go install ariga.io/atlas/cmd/atlas@latest

tooling-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Fixture

fixture-dahua-push:
	curl -s -H "Content-Type: application/json" --data-binary @fixtures/dahua.json localhost:8080/v1/dahua | jq

fixture-dahua-list:
	jq -r 'keys | join("\n")' fixtures/dahua.json
