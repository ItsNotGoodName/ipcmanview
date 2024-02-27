export DIR=ipcmanview_data
export VITE_HOST=127.0.0.1
export WEB_PATH=internal/web

export CMD_AIR=github.com/cosmtrek/air@v1.49.0
export CMD_TASK=github.com/go-task/task/v3/cmd/task@v3.34.1
export CMD_GOOSE=github.com/pressly/goose/v3/cmd/goose@v3.18.0
export CMD_SQLC=github.com/sqlc-dev/sqlc/cmd/sqlc@v1.25.0
export CMD_TWIRP=github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.3
export CMD_PROTOC_GEN_GO=google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0

export PROTOC_VERSION=25.1
export PROTOC_ZIP=protoc-$(PROTOC_VERSION)-linux-x86_64.zip

-include .env

_:
	mkdir -p $(WEB_PATH)/dist $(DIR) && touch $(WEB_PATH)/dist/index.html

clean:
	rm -rf $(DIR)

migrate:
	goose -dir internal/sqlite/migrations sqlite3 $(DIR)/sqlite.db up

generate:
	go generate ./...

run:
	go run ./cmd/ipcmanview serve

debug:
	go run ./cmd/ipcmanview debug

preview: generate run

migration:
	atlas migrate diff $(name) --env local

hash:
	atlas migrate hash --env local

# ---------- Dev

# Start backend
dev:
	go run $(CMD_AIR)

# Start frontend
dev-web:
	cd $(WEB_PATH) && pnpm install && pnpm run dev

# Proxy request to frontend and backend
dev-proxy:
	go run ./scripts/dev-proxy

# ---------- Gen

# Generate code
gen: gen-sqlc gen-pubsub gen-bus gen-proto gen-typescriptify

gen-sqlc:
	go run $(CMD_SQLC) generate

gen-pubsub:
	sh ./scripts/generate-pubsub-events.sh ./internal/event/models.go

gen-bus:
	go run ./scripts/generate-bus ./internal/event/models.go

gen-proto:
	cd rpc && protoc --go_out=. --twirp_out=. rpc.proto
	cd $(WEB_PATH) && mkdir -p ./src/twirp && pnpm exec protoc --ts_out=./src/twirp --ts_opt=generate_dependencies --proto_path=$(shell readlink -f rpc) rpc.proto

gen-typescriptify:
	go run ./scripts/typescriptify $(WEB_PATH)/src/lib/models.gen.ts

# ---------- Tooling

# Install tooling
tooling: tooling-air tooling-task tooling-goose tooling-sqlc tooling-twirp tooling-protoc-gen-go tooling-protoc-gen-ts

tooling-air:
	go install $(CMD_AIR)

tooling-task:
	go install $(CMD_TASK)

tooling-goose:
	go install $(CMD_GOOSE)

tooling-sqlc:
	go install $(CMD_SQLC)

tooling-twirp:
	go install $(CMD_TWIRP)

tooling-protoc-gen-go:
	go install $(CMD_PROTOC_GEN_GO)

tooling-protoc-gen-ts:
	cd $(WEB_PATH) && pnpm install

# ---------- Install

install-protoc:
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)
	sudo unzip -o $(PROTOC_ZIP) -d /usr/local bin/protoc
	sudo unzip -o $(PROTOC_ZIP) -d /usr/local 'include/*'
	rm $(PROTOC_ZIP)

install-atlas:
	# TODO: pin atlas version
	curl -OL https://release.ariga.io/atlas/atlas-community-linux-amd64-latest
	chmod +x atlas-community-linux-amd64-latest
	mv atlas-community-linux-amd64-latest ~/.local/bin/atlas

# Workflow

workflow-tooling: tooling-task tooling-sqlc tooling-twirp tooling-protoc-gen-go tooling-protoc-gen-ts

workflow-nightly:
	go run $(CMD_TASK) nightly
