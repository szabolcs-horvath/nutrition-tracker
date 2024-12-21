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
  - Install them with:
    ```shell
    make go-deps
    ```

### Set up
1. #### Create the database
    ```shell
    sqlite3 sqlite/nutrition-tracker.db
    ```
    You can exit from the interactive console with `.quit`

2. #### Build the project
    ```shell
    make build
    ```

## Running
```shell
out/nutrition-tracker
```

## Adding a new migration
```shell
make create-migration MIGRATION_NAME=<migration_name>
```