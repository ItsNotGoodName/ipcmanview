gen:
	jet -dsn=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable -path=./internal/dbgen
	jet -dsn=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable -path=./internal/dbgen -schema dahua
	webrpc-gen -schema=./server/api.ridl -target=golang -pkg=service -server -out=./server/service/proto.gen.go
	webrpc-gen -schema=./server/api.ridl -target=typescript -client -out=./ui/src/core/client.gen.ts

cli:
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres" go run ./cmd/console

debug:
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres" go run ./cmd/debug

preview: build-ui server

build-ui:
	cd ui && pnpm run build && cd ..

server:
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres" go run ./cmd/server

migrate:
	cd migrations && tern migrate

dev:
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres" air

dev-db:
	podman run --rm -e POSTGRES_PASSWORD=postgres -p 5432:5432 docker.io/postgres:15 -c log_statement=all

dev-ui:
	cd ui && pnpm run dev

dep: dep-tern dep-jet dep-air dep-webrpc-gen dep-ui

dep-tern:
	go install github.com/jackc/tern/v2@latest

dep-jet:
	go install github.com/go-jet/jet/v2/cmd/jet@latest

dep-air:
		go install github.com/cosmtrek/air@latest

# TODO: fix this install
dep-webrpc-gen:
		go install -ldflags="-s -w -X github.com/webrpc/webrpc.VERSION=v0.12.0" github.com/webrpc/webrpc/cmd/webrpc-gen@v0.12.0

dep-ui:
	cd ui && pnpm install
