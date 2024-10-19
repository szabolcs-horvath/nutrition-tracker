MIGRATE_CLI_LOCATION ?= vendor/migrate-cli
MIGRATE_CLI_BINARY_LOCATION ?= $(MIGRATE_CLI_LOCATION)/cmd/migrate
MIGRATE_CLI_DATABASE_DRIVERS ?= sqlite3
MIGRATE_CLI_SOURCES ?= file
MIGRATE_VERSION ?= v4.18.1
GOOS ?= darwin
GOARCH ?= amd64
SQLITE_DB_FILE ?= sqlite/nutrition-tracker.db
SQLITE_MIGRATIONS_DIR ?= sqlite/migrations

build:
	sqlc generate
	go build -o out/nutrition-tracker

clean:
	rm -rf /generated /out

migrate-up: build-migrate-cli
	$(MIGRATE_CLI_BINARY_LOCATION)/migrate -source file://$(SQLITE_MIGRATIONS_DIR) -database sqlite3://$(SQLITE_DB_FILE) up

build-migrate-cli:
ifeq ("$(wildcard $(MIGRATE_CLI_LOCATION)/cmd/migrate/migrate)", "")
	@echo "The migrate-cli location exists but the binary doesn't"
	@echo "Building the migrate-cli with the needed database drivers and migration file sources"
	make clean-migrate-cli

	git clone git@github.com:golang-migrate/migrate.git $(MIGRATE_CLI_LOCATION)
	git -C $(MIGRATE_CLI_LOCATION) fetch --all
	git -C $(MIGRATE_CLI_LOCATION) checkout tags/$(MIGRATE_VERSION)

	cd $(MIGRATE_CLI_BINARY_LOCATION) && CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -o migrate -ldflags='-X main.Version=$(VERSION)' -tags '$(MIGRATE_CLI_DATABASE_DRIVERS) $(MIGRATE_CLI_SOURCES)' .
endif

clean-migrate-cli:
	rm -rf $(MIGRATE_CLI_LOCATION)

