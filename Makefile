# migration name
NAME ?= name
TEST_DB_PASSWORD ?= password
TEST_DB_HOST_PORT ?= 5432
TEST_DB_URL ?= postgres://postgres:$(TEST_DB_PASSWORD)@localhost:$(TEST_DB_HOST_PORT)/postgres?sslmode=disable 

.PHONE: query

dev:
	@VERSION=dev APP_ENV=development air
up:
	@goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
down:
	@goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down-to 0
migration:
	@goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" create $(NAME) sql

init-dev-db:
	@docker run --name konbini-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DATABASE=konbini -p 5432:5432 -d postgres:14 && \
		timeout 90s bash -c "until docker exec konbini-postgres pg_isready ; do sleep 5 ; done" && echo "Postgres is ready! Run migrations with 'make up'"
start-dev-db:
	@docker start konbini-postgres && \
		timeout 90s bash -c \
		"until docker exec konbini-postgres pg_isready ; do sleep 5 ; done" && \
		echo "Postgres is ready! Run migrations with 'make up'"
stop-dev-db:
	@docker stop konbini-postgres
check-dev-db:
	@docker exec konbini-postgres pg_isready
clean-dev-db:
	@docker stop konbini-postgres && docker rm konbini-postgres
query:
	@docker exec konbini-postgres psql -U postgres -d postgres -c "$(filter-out $@,$(MAKECMDGOALS))"
# Prevent make from treating the query string as a target
%:
	@:

test: start-testdb
	@trap 'make stop-testdb' EXIT; \
		if ! make run-tests; then \
			exit 1; \
		fi

start-testdb:
	@echo "Setting up PostgreSQL database..."
	@docker run --name test-postgres -e POSTGRES_PASSWORD=$(TEST_DB_PASSWORD) -d -p "$(TEST_DB_HOST_PORT):5432" postgres:14
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker exec test-postgres pg_isready -U postgres; do sleep 1; done
	@goose -dir ./migrations postgres $(TEST_DB_URL) up

run-tests:
	APP_ENV=test VERSION=0.0.0-test PASS_ENCRYPT_ALGO=md5 DB_URL=$(TEST_DB_URL) go test -v ./router ./store

stop-testdb:
	@docker container stop test-postgres
	@docker container rm -f test-postgres
