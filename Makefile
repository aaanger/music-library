build:
	docker compose build
run:
	docker compose up
migrate:
	goose -dir pkg/db/migrations postgres "host=${PSQL_HOST} port=${PSQL_PORT} user=${PSQL_USERNAME} password=${PSQL_PASSWORD} dbname=${PSQL_DBNAME} sslmode=${PSQL_SSLMODE}" up
rollback:
	goose -dir pkg/db/migrations postgres "host=${PSQL_HOST} port=${PSQL_PORT} user=${PSQL_USERNAME} password=${PSQL_PASSWORD} dbname=${PSQL_DBNAME} sslmode=${PSQL_SSLMODE}" down