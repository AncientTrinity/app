## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo '--Running application--'
	docker compose run --rm go go run ./cmd/api -port=8081 -env=development -db-dsn=$${COMMENTS_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo "Creating migration files for ${name}..."
	docker run --rm -v $$(pwd)/migrations:/migrations migrate/migrate create -seq -ext=.sql -dir=/migrations ${name}

## db/migrations/up: apply all migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo "Running up migrations..."
	docker compose run --rm migrate up

## db/migrations/down: rollback last migration
.PHONY: db/migrations/down
db/migrations/down:
	@echo "Rolling back last migration..."
	docker compose run --rm migrate down 1

## db/migrations/force version=$1: force database version
.PHONY: db/migrations/force
db/migrations/force:
	@echo "Forcing version to ${version}..."
	docker compose run --rm migrate force ${version}

## db/psql: connect to PostgreSQL using psql
.PHONY: db/psql
db/psql:
	docker exec -it postgres-db psql -U user -d mydb
