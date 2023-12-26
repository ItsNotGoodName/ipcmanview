export DIR=ipcmanview_data
export VITE_HOST=127.0.0.1

-include .env

_:
	mkdir web/dist -p && touch web/dist/index.html

migrate:
	goose -dir internal/migrations/sql sqlite3 "$(DIR)/sqlite.db" up

clean:
	rm -rf $(DIR)

generate:
	go generate ./...

run:
	go run ./cmd/ipcmanview serve

preview: generate run

# Dev

dev:
	air

dev-web:
	cd internal/web && pnpm install && pnpm run dev

dev-webnext:
	cd internal/webnext && pnpm install && pnpm run dev

# Gen

gen: gen-sqlc gen-pubsub gen-bus

gen-sqlc:
	sqlc generate

gen-pubsub:
	sh ./scripts/generate-pubsub-events.sh ./internal/models/event.go

gen-bus:
	go run ./scripts/generate-bus.go -input ./internal/models/event.go -output ./internal/core/bus.gen.go

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
