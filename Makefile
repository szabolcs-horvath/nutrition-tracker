SQLITE_DB_FILE ?= sqlite/nutrition-tracker.db
IT_SQLITE_DB_FILE ?= it/nutrition-tracker-test.db
SQLITE_MIGRATIONS_DIR ?= sqlite/migrations
SQLC_VERSION ?= v1.27.0
GOLANG_MIGRATE_VERSION ?= v4.18.1
STRINGER_VERSION ?= v0.28.0
HTMX_VERSION ?= 2.0.3
BOOTSTRAP_VERSION ?= 5.3.3
GOCOVERDIR ?= coverage
CGO_ENABLED=1 # Required for sqlite3 driver

install-go-deps:
	go install -v github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	go install -v -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@$(GOLANG_MIGRATE_VERSION)
	go install -v golang.org/x/tools/cmd/stringer@$(STRINGER_VERSION)

init-db: migrate-up
	sqlite3 $(SQLITE_DB_FILE) < sqlite/seed.sql

build: sqlc generate
	go build -o out/nutrition-tracker -mod=readonly

sqlc:
	sqlc generate

generate:
	go generate ./...

clean:
	rm -rf generated out

ut: sqlc
	rm -f $(GOCOVERDIR)/ut-coverage.out
	mkdir -p $(GOCOVERDIR)
	go test ./... -v -coverprofile=$(GOCOVERDIR)/ut-coverage.out -covermode=atomic -coverpkg=./...

it: sqlc init-test-db
	rm -rf $(GOCOVERDIR)/it-coverage
	mkdir -p $(GOCOVERDIR)/it-coverage
	go build -o out/nutrition-tracker-it -mod=readonly -covermode=atomic
	it/integration-test.sh out/nutrition-tracker-it $(GOCOVERDIR)/it-coverage
	go tool covdata textfmt -i=$(GOCOVERDIR)/it-coverage -o=$(GOCOVERDIR)/it-coverage.out

init-test-db:
	rm -f $(IT_SQLITE_DB_FILE)
	make migrate-up-file SPECIFIED_DB_FILE=$(IT_SQLITE_DB_FILE)
	sqlite3 $(IT_SQLITE_DB_FILE) < it/seed_test.sql

coverage: ut it
	go run ./.github/merge-coverprofiles.go $(GOCOVERDIR)/merged-coverage.out $(GOCOVERDIR)/ut-coverage.out $(GOCOVERDIR)/it-coverage.out
	go tool cover -html=$(GOCOVERDIR)/merged-coverage.out

coverage-ut: ut
	go tool cover -html=$(GOCOVERDIR)/ut-coverage.out

coverage-it: it
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
