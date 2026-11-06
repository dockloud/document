-include  .env
export DB_URL := $(DATABASE_URL)

goose_up:
	goose -dir internal/infrastructure/database/postgres/migrations \
	postgres "$$DB_URL" up