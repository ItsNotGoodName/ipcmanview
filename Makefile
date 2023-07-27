start:
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres" go run .

dev-db:
	podman run --rm -e POSTGRES_PASSWORD=postgres -p 5432:5432 docker.io/postgres:15

dep-tern:
	go install github.com/jackc/tern/v2@latest
