start:
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres" go run .

gen:
	jet -dsn=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable -path=./internal/db/gen

dev-db:
	podman run --rm -e POSTGRES_PASSWORD=postgres -p 5432:5432 docker.io/postgres:15

dep: dep-tern dep-jet

dep-tern:
	go install github.com/jackc/tern/v2@latest

dep-jet:
	go install github.com/go-jet/jet/v2/cmd/jet@latest
