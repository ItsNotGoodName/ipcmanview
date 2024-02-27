export DIR=ipcmanview_data
export VITE_HOST=127.0.0.1
export WEB_PATH=internal/web

-include .env

_:
	mkdir -p "$(WEB_PATH)/dist" "$(DIR)" && touch "$(WEB_PATH)/dist/index.html"

clean:
	rm -rf $(DIR)

migrate:
	goose -dir internal/sqlite/migrations sqlite3 "$(DIR)/sqlite.db" up

build:
	go generate ./...

run:
	go run ./cmd/ipcmanview serve

debug:
	go run ./cmd/ipcmanview debug

preview: build run

nightly:
	task nightly

migration:
	atlas migrate diff $(name) --env local

hash:
	atlas migrate hash --env local

# Dev

dev:
	air

dev-proxy:
	go run ./scripts/dev-proxy

dev-web:
	cd "$(WEB_PATH)" && pnpm install && pnpm run dev

# Gen

gen: gen-sqlc gen-pubsub gen-bus gen-proto gen-typescriptify

gen-sqlc:
	sqlc generate

gen-pubsub:
	sh ./scripts/generate-pubsub-events.sh ./internal/event/models.go

gen-bus:
	go run ./scripts/generate-bus ./internal/event/models.go

gen-proto:
	cd rpc && protoc --go_out=. --twirp_out=. rpc.proto
	cd "$(WEB_PATH)" && mkdir -p ./src/twirp && pnpm exec protoc --ts_out=./src/twirp --ts_opt=generate_dependencies --proto_path=../../rpc rpc.proto

gen-typescriptify:
	go run ./scripts/typescriptify ./internal/web/src/lib/models.gen.ts

# Tooling

build-tooling: tooling-task tooling-sqlc tooling-twirp tooling-protoc-gen-go tooling-protoc-gen-ts

tooling: tooling-air tooling-task tooling-goose tooling-atlas tooling-sqlc tooling-twirp tooling-protoc-gen-go tooling-protoc-gen-ts

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

tooling-twirp:
	go install github.com/twitchtv/twirp/protoc-gen-twirp@latest

tooling-protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

tooling-protoc-gen-ts:
	cd "$(WEB_PATH)" && pnpm install

# Install

install-protoc:
	PROTOC_ZIP=protoc-25.1-linux-x86_64.zip
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v25.1/$PROTOC_ZIP
	sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
	sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
