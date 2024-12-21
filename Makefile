SQLITE_DB_FILE ?= sqlite/nutrition-tracker.db
SQLITE_MIGRATIONS_DIR ?= sqlite/migrations
SQLC_VERSION ?= v1.27.0
GOLANG_MIGRATE_VERSION ?= v4.18.1
HTMX_VERSION ?= 2.0.3
BOOTSTRAP_VERSION ?= 5.3.3
GOCOVERDIR ?= coverage

go-deps:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@$(GOLANG_MIGRATE_VERSION)

init-db: migrate-up
	sqlite3 $(SQLITE_DB_FILE) < sqlite/seed.sql

build: sqlc
	go build -o out/nutrition-tracker -mod=readonly

sqlc:
	sqlc generate

clean:
	rm -rf generated out

unit-test: sqlc
	rm -f $(GOCOVERDIR)/unit-coverage.out
	mkdir -p $(GOCOVERDIR)
	go test ./... -v -coverprofile=$(GOCOVERDIR)/unit-coverage.out -covermode=atomic -coverpkg=./...

integration-test: sqlc
	rm -rf $(GOCOVERDIR)/it-coverage
	mkdir -p $(GOCOVERDIR)/it-coverage
	go build -o out/nutrition-tracker-integration-test -mod=readonly -covermode=atomic
	integration-test/integration-test.sh out/nutrition-tracker-integration-test $(GOCOVERDIR)/it-coverage
	go tool covdata textfmt -i=$(GOCOVERDIR)/it-coverage -o=$(GOCOVERDIR)/it-coverage.out

coverage: unit-test integration-test
	go run ./.github/merge-coverprofiles.go $(GOCOVERDIR)/merged-coverage.out $(GOCOVERDIR)/unit-coverage.out $(GOCOVERDIR)/it-coverage.out
	go tool cover -html=$(GOCOVERDIR)/merged-coverage.out

coverage-unit: unit-test
	go tool cover -html=$(GOCOVERDIR)/unit-coverage.out

coverage-it: integration-test
	go tool cover -html=$(GOCOVERDIR)/it-coverage.out

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
