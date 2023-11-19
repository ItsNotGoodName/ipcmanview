-include .env

dev-gateway:
	air

tooling: tooling-air tooling-task

# Tooling

tooling-air:
	go install github.com/cosmtrek/air@latest

tooling-task:
	go install github.com/go-task/task/v3/cmd/task@latest

# Fixture

fixture-dahua-push:
	curl -s -H "Content-Type: application/json" --data-binary @fixtures/dahua.json localhost:8080/v1/dahua | jq

fixture-dahua-list:
	jq -r 'keys | join("\n")' fixtures/dahua.json
