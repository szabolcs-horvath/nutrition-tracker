SQLITE_DB_FILE ?= sqlite/nutrition-tracker.db
SQLITE_MIGRATIONS_DIR ?= sqlite/migrations
HTMX_VERSION ?= 2.0.3
BOOTSTRAP_VERSION ?= 5.3.3

build: migrate-up sqlc
	go build -o out/nutrition-tracker -mod=readonly

sqlc:
	sqlc generate

clean:
	rm -rf generated out

create-migration:
ifneq ($(MIGRATION_NAME),)
	migrate create -dir sqlite/migrations -ext .sql $(MIGRATION_NAME)
else
	@echo "MIGRATION_NAME needs to be specified. Usage: make create-migration MIGRATION_NAME=<migration_name>"
	exit 1
endif

migrate-up:
	migrate -source file://$(SQLITE_MIGRATIONS_DIR) -database sqlite3://$(SQLITE_DB_FILE) up

migrate-up-1:
	migrate -source file://$(SQLITE_MIGRATIONS_DIR) -database sqlite3://$(SQLITE_DB_FILE) up 1

migrate-up-file:
ifneq ($(SPECIFIED_DB_FILE),)
	migrate -source file://$(SQLITE_MIGRATIONS_DIR) -database sqlite3://$(SPECIFIED_DB_FILE) up
else
	@echo "SPECIFIED_DB_FILE needs to be specified. Usage: make migrate-up-file SPECIFIED_DB_FILE=<db_file>"
	exit 1
endif

migrate-up-1-file:
ifneq ($(SPECIFIED_DB_FILE),)
	migrate -source file://$(SQLITE_MIGRATIONS_DIR) -database sqlite3://$(SPECIFIED_DB_FILE) up 1
else
	@echo "SPECIFIED_DB_FILE needs to be specified. Usage: make migrate-up-1-file SPECIFIED_DB_FILE=<db_file>"
	exit 1
endif

migrate-down-1:
	migrate -source file://$(SQLITE_MIGRATIONS_DIR) -database sqlite3://$(SQLITE_DB_FILE) down 1

migrate-down-1-file:
ifneq ($(SPECIFIED_DB_FILE),)
	migrate -source file://$(SQLITE_MIGRATIONS_DIR) -database sqlite3://$(SPECIFIED_DB_FILE) down 1
else
	@echo "SPECIFIED_DB_FILE needs to be specified. Usage: make migrate-down-1-file SPECIFIED_DB_FILE=<db_file>"
	exit 1
endif


download-htmx:
	cd ./web/static/vendor/htmx; \
	curl -O https://unpkg.com/htmx.org@$(HTMX_VERSION)/dist/htmx.min.js; \
	cd -;

download-bootstrap:
	cd ./web/static/vendor/bootstrap; \
	curl -O https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/css/bootstrap.min.css; \
	curl -O https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/js/bootstrap.bundle.min.js; \
	cd -;
