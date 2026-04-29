include .env
export

CONNECTION_STR = postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

DEV_COMPOSE=docker-compose.dev.yml
ENV_FILE=.env
access-db:
	docker exec -it postgres-db psql -U oxnen -d web-socket-app
server:
	cd cmd/api && go run .
sqlc:
	sqlc generate
migrate-create:
	migrate create -ext sql -dir=$(pwd)./internal/database/migrations -seq $(NAME)
migrate-up:
	migrate -path=$(pwd)./internal/database/migrations -database "$(CONNECTION_STR)" up
migrate-down:
	migrate -path=$(pwd)./internal/database/migrations -database "$(CONNECTION_STR)" down 1
migrate-force:
	migrate -path=$(pwd)./internal/database/migrations -database "$(CONNECTION_STR)" force $(VERSION)
migrate-goto:
	migrate -path=$(pwd)./internal/database/migrations -database "$(CONNECTION_STR)" goto $(VERSION)
migrate-check-version:
	migrate -path=$(pwd)./internal/database/migrations -database "$(CONNECTION_STR)" version
dev:
	docker compose -f $(DEV_COMPOSE) down
	docker compose -f $(DEV_COMPOSE) --env-file $(ENV_FILE) up --build
stop-dev:
	docker compose -f $(DEV_COMPOSE) down
.PHONY:stop-dev server sqlc migrate-create migrate-up migrate-down migrate-goto migrate-force migrate-check-version dev access-db
