export DIR=ipcmanview_data
export VITE_HOST=127.0.0.1
export WEB_PATH=internal/web

-include .env

TOOL_AIR=github.com/cosmtrek/air@v1.51.0
TOOL_TASK=github.com/go-task/task/v3/cmd/task@v3.35.1
TOOL_GOOSE=github.com/pressly/goose/v3/cmd/goose@v3.18.0
TOOL_SQLC=github.com/sqlc-dev/sqlc/cmd/sqlc@v1.25.0
TOOL_TWIRP=github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.3
TOOL_PROTOC_GEN_GO=google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
TOOL_OAPI_CODEGEN=github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0

PROTOC_VERSION=25.1
PROTOC_ZIP=protoc-$(PROTOC_VERSION)-linux-x86_64.zip

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
	air

# Start frontend
dev-web:
	cd $(WEB_PATH) && pnpm install && pnpm run dev

# Proxy request to frontend and backend
dev-proxy:
	go run ./scripts/dev-proxy

# ---------- Gen

# Generate code
gen: gen-sqlc gen-proto gen-mediamtx gen-pubsub gen-hub gen-typescriptify 

gen-sqlc:
	sqlc generate

gen-proto:
	cd rpc && protoc --go_out=. --twirp_out=. rpc.proto
	cd $(WEB_PATH) && mkdir -p ./src/twirp && pnpm exec protoc --ts_out=./src/twirp --ts_opt=generate_dependencies --proto_path=$(shell readlink -f rpc) --ts_opt long_type_string rpc.proto

gen-mediamtx:
	oapi-codegen -package mediamtx ./internal/mediamtx/swagger.json > ./internal/mediamtx/mediamtx.gen.go

gen-pubsub:
	sh ./scripts/generate-pubsub-events.sh ./internal/bus/models.go

gen-hub:
	go run ./scripts/generate-hub ./internal/bus/models.go

gen-typescriptify:
	go run ./scripts/typescriptify $(WEB_PATH)/src/lib/models.gen.ts

# ---------- Tooling

# Install tooling
tooling: tooling-air tooling-task tooling-goose tooling-sqlc tooling-twirp tooling-protoc-gen-go tooling-protoc-gen-ts tooling-oapi-codegen

tooling-build: tooling-task tooling-sqlc tooling-twirp tooling-protoc-gen-go tooling-protoc-gen-ts tooling-oapi-codegen

tooling-air:
	go install $(TOOL_AIR)

tooling-task:
	go install $(TOOL_TASK)

tooling-goose:
	go install $(TOOL_GOOSE)

tooling-sqlc:
	go install $(TOOL_SQLC)

tooling-twirp:
	go install $(TOOL_TWIRP)

tooling-protoc-gen-go:
	go install $(TOOL_PROTOC_GEN_GO)

tooling-protoc-gen-ts:
	cd $(WEB_PATH) && pnpm install

tooling-oapi-codegen:
	go install $(TOOL_OAPI_CODEGEN)

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
