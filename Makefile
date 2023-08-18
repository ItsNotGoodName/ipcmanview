export DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres
export JWT_SECRET=stop-logging-me-out

-include .env

# NOTE: IDK if the wildcard is doing anything
.PHONY: fake debug server dev-* dep-*

gen: gen-jet gen-webrpc

gen-jet:
	jet -dsn=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable -path=./internal/dbgen
	jet -dsn=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable -path=./internal/dbgen -schema dahua

gen-webrpc:
	webrpc-gen -schema=./server/api.ridl -target=golang -pkg=rpcgen -server -out=./server/rpcgen/rpcgen.gen.go
	webrpc-gen -schema=./server/api.ridl -target=./server/gen-typescript-nuxt -client -out=./ui/core/client.gen.ts

preview: build-ui server

preview-ui: build-ui
	cd ui && npm run preview

fake:
	go run ./cmd/ipcmanview-fake

debug:
	go run ./cmd/ipcmanview-debug

server:
	go run ./cmd/ipcmanview

build-ui:
	cd ui && pnpm run build && cd ..

dev-db:
	podman run --rm -e POSTGRES_PASSWORD=postgres -p 5432:5432 docker.io/postgres:15 -c log_statement=all

dev-fake:
	air -build.cmd="go build -o ./tmp/main -tags dev ./cmd/ipcmanview-fake"

dev-debug:
	air -build.cmd="go build -o ./tmp/main -tags dev ./cmd/ipcmanview-debug"

dev-server:
	air

dev-ui:
	cd ui && pnpm run dev

dev-migrate:
	cd migrations && tern migrate

# Development dependencies

dep: dep-tern dep-jet dep-air dep-webrpc-gen dep-ui

dep-tern:
	go install github.com/jackc/tern/v2@latest

dep-jet:
	go install github.com/go-jet/jet/v2/cmd/jet@latest

dep-air:
		go install github.com/cosmtrek/air@latest

dep-webrpc-gen:
		go install -ldflags="-s -w -X github.com/webrpc/webrpc.VERSION=v0.12.1" github.com/webrpc/webrpc/cmd/webrpc-gen@v0.12.1

dep-ui:
	cd ui && pnpm install
