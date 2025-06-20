include .env
export

MIGRATE=migrate
MIGRATIONS_DIR=./migrations
DB_URL=mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp(db:$(MYSQL_PORT))/$(MYSQL_NAME)?multiStatements=true

migrate-new:
	@$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-goto:
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" goto $(version)

migrate-version:
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version
