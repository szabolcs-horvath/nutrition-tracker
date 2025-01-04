SQLITE_DB_FILE ?= sqlite/nutrition-tracker.db
IT_SQLITE_DB_FILE ?= it/nutrition-tracker-test.db
SQLITE_MIGRATIONS_DIR ?= sqlite/migrations
GOCOVERDIR ?= coverage

SQLC_VERSION ?= v1.27.0
GOLANG_MIGRATE_VERSION ?= v4.18.1
STRINGER_VERSION ?= v0.28.0
DELVE_VERSION ?= latest
HTMX_VERSION ?= 2.0.3
BOOTSTRAP_VERSION ?= 5.3.3
BOOTSTRAP_ICONS_VERSION ?= 1.11.3

install: install-go-deps init-db build
	cp -u prod/nutrition-tracker.service /etc/systemd/system/nutrition-tracker.service
	systemctl enable nutrition-tracker.service

start:
	systemctl start nutrition-tracker.service

install-go-deps:
	go install -v github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	CGO_ENABLED=1 go install -v -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@$(GOLANG_MIGRATE_VERSION)
	go install -v golang.org/x/tools/cmd/stringer@$(STRINGER_VERSION)

init-db: migrate-up
	sqlite3 $(SQLITE_DB_FILE) < sqlite/seed.sql

build: sqlc generate
	go build -o out/nutrition-tracker -mod=readonly

debug: sqlc generate
	go install -v github.com/go-delve/delve/cmd/dlv@$(DELVE_VERSION)
	go build -o out/nutrition-tracker-debug -mod=readonly -gcflags="all=-N -l"
	dlv --listen=:443 --headless=true --api-version=2 --accept-multiclient exec ./out/nutrition-tracker-debug

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

download-bootstrap-icons:
	cd ./web/static/vendor/bootstrap/icons; \
	rm -rf ./*; \
	curl -LO https://github.com/twbs/icons/releases/download/v$(BOOTSTRAP_ICONS_VERSION)/bootstrap-icons-$(BOOTSTRAP_ICONS_VERSION).zip; \
	unzip -d . bootstrap-icons-$(BOOTSTRAP_ICONS_VERSION).zip; \
	mv -f bootstrap-icons-$(BOOTSTRAP_ICONS_VERSION)/* .; \
	rm -rf bootstrap-icons-$(BOOTSTRAP_ICONS_VERSION); \
	rm bootstrap-icons-$(BOOTSTRAP_ICONS_VERSION).zip; \
	cd -;

get-cert-from-letsencrypt-interactive:
	snap install --classic certbot
	ln -s /snap/bin/certbot /usr/bin/certbot
	certbot certonly --standalone

create-self-signed-cert:
	mkdir -p certs
	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
		-out certs/cert.pem \
		-keyout certs/key.pem \
		-subj "/CN=localhost"

create-self-signed-cert-interactive:
	mkdir -p certs
	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
		-out certs/cert.pem \
		-keyout certs/key.pem
