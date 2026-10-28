SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -ExecutionPolicy Bypass -Command
.PHONY: up down run build deps swagger migrate-up migrate-down migrate-version logs test
up:
	docker compose up -d todoapp-postgres port-forwarder
	docker compose run --rm todoapp-postgres-migrate up
down:
	docker compose down
run:
	./scripts/run.ps1
build:
	go build -o out/kontakt.exe ./cmd/app
deps:
	go mod download
swagger:
	docker compose run --rm swagger
migrate-up:
	docker compose run --rm todoapp-postgres-migrate up
migrate-down:
	docker compose run --rm todoapp-postgres-migrate down 1
migrate-version:
	docker compose run --rm todoapp-postgres-migrate version
logs:
	Get-Content out/logs/app.log -Tail 50 -Wait
test:
	go test ./...
