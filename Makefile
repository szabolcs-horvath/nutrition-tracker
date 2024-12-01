SQLITE_DB_FILE ?= sqlite/nutrition-tracker.db
SQLITE_MIGRATIONS_DIR ?= sqlite/migrations
HTMX_VERSION ?= 2.0.3
BOOTSTRAP_VERSION ?= 5.3.3

build: migrate-up sqlc
	go build -o out/nutrition-tracker -mod=readonly

sqlc:
	sqlc generate

clean:
	rm -rf generated out vendor

create-migration:
ifneq ($(MIGRATION_NAME),)
	migrate create -dir sqlite/migrations -ext .sql $(MIGRATION_NAME)
else
	@echo "MIGRATION_NAME needs to be specified. Usage: make create-migration MIGRATION_NAME=<migration_name>"
	exit 1
endif

migrate-up:
	migrate -source file://$(SQLITE_MIGRATIONS_DIR) -database sqlite3://$(SQLITE_DB_FILE) up

download-htmx:
	cd ./static/vendor/htmx; \
	curl -O https://unpkg.com/htmx.org@$(HTMX_VERSION)/dist/htmx.min.js; \
	cd -;

download-bootstrap:
	cd ./static/vendor/bootstrap; \
	curl -O https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/css/bootstrap.min.css; \
	curl -O https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/js/bootstrap.bundle.min.js; \
	cd -; \
