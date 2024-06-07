# migration name
NAME ?= name

dev:
	@CI=true APP_VERSION=dev APP_ENV=development air
up:
	@goose -dir ./migrations postgres "postgres://konbini:konbini@localhost:5432/konbini?sslmode=disable" up
down:
	@goose -dir ./migrations postgres "postgres://konbini:konbini@localhost:5432/konbini?sslmode=disable" down-to 0
migration:
	@goose -dir ./migrations postgres "postgres://konbini:konbini@localhost:5432/konbini?sslmode=disable" create $(NAME) sql
