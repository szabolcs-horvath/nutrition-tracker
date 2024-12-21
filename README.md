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
  - :white_check_mark: [`sqlc`](https://docs.sqlc.dev/en/latest/overview/install.html)
  - :white_check_mark: [`golang-migrate`](https://github.com/golang-migrate/migrate/tree/master)
  - Install these with:
    ```shell
    make go-deps
    ```

### Set up
1. #### Initialize the database
    ```shell
    make init-db
    ```

2. #### Build the project
    ```shell
    make build
    ```

## Running
```shell
out/nutrition-tracker <path-to-.env-file>
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

## Adding a new migration
```shell
make create-migration MIGRATION_NAME=<migration_name>
```