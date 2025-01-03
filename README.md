![coverage](https://raw.githubusercontent.com/szabolcs-horvath/nutrition-tracker/badges/.badges/main/coverage.svg)

# Nutrition Tracker for Kinga ( :shushing_face: SECRET :shushing_face: )
This project's purpose is to help Kinga track and plan her diet.

## Features
- N/A

## Building
### Prerequisites
- :white_check_mark: [`go@1.23.1`](https://go.dev/dl/)
- :white_check_mark: [`sqlite3`](https://command-not-found.com/sqlite3)
- :white_check_mark: go dependencies
  - :white_check_mark: [`sqlc`](https://docs.sqlc.dev/en/latest/)
  - :white_check_mark: [`golang-migrate`](https://github.com/golang-migrate/migrate/)
  - :white_check_mark: [`stringer`](https://pkg.go.dev/golang.org/x/tools/cmd/stringer)
  - Install these with:
    ```shell
    make go-deps
    ```

### Set up
1. #### Create a `.env` file with the necessary environment variables
    ```shell
    DB_FILE="sqlite/nutrition-tracker.db"
    PORT="6969" # Optional, defaults to "80" or "443" depending on whether TLS is enabled
    
    TLS_DISABLED="false" # Optional, defaults to "false"
    TLS_CERT_FILE="certs/cert.pem" # Optional, defaults to "certs/cert.pem"
    TLS_KEY_FILE="certs/key.pem" # Optional, defaults to "certs/key.pem"
   
    AUTH0_DISABLED="false" # Optional, defaults to "false"
    AUTH0_DOMAIN="<AUTH0_DOMAIN>" # Can be found in the Auth0 dashboard
    AUTH0_CLIENT_ID="<AUTH0_CLIENT_ID>" # Can be found in the Auth0 dashboard
    AUTH0_CLIENT_SECRET="<AUTH0_CLIENT_SECRET>" # Can be found in the Auth0 dashboard
    AUTH0_CALLBACK_URL="https://nutrition-tracking.com/auth/callback" # Should be the same as the one set in the Auth0 dashboard
    COOKIE_STORE_AUTH_KEY="<KEY_TO_AUTHENTICATE_THE_COOKIE_STORE>" # Should ideally be at least 64 bytes of random data
    ```
2. #### Initialize the database
    ```shell
    make init-db
    ```
3. #### Build the project
    ```shell
    make build
    ```

## Running
###### (The .env file only needs to be specified if it is not in the project root)
```shell
out/nutrition-tracker <PATH_TO_YOUR_ENV_FILE>
```

## Checking test coverage
###### Both unit and integration tests
```shell
make coverage
```
###### Only unit tests
```shell
make unit-coverage
```
###### Only integration tests
```shell
make integration-coverage
```

## Migrations
###### Create a new migration
```shell
make create-migration MIGRATION_NAME=<MIGRATION_NAME>
```
###### Apply all migrations
```shell
make migrate-up
```
###### Rollback the last migration
```shell
make migrate-down-1
```