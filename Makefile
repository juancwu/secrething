# migration name
NAME ?= name
DB_URL ?= http://localhost:8080

dev:
	@VERSION=dev APP_ENV=development air
up:
	DATABASE_URL=$(DB_URL) geni up
down:
	DATABASE_URL=$(DB_URL) geni down
migration:
	DATABASE_URL=$(DB_URL) geni new $(NAME)
