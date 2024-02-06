.PHONY: run

create-db:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -d postgres:15-alpine

create-table:
	docker exec -it postgres15 createdb --username=root --owner=root chat

migrate:
	migrate -path db/migrations -database "postgresql://root:123@localhost:5433/chat?sslmode=disable" -verbose up

create-redis:
	docker run --name chat-redis -p 6379:6379 -d redis:latest

generate-api-docs:
	swag init -g cmd/server/main.go --output docs
