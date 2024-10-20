# Nutrition Tracker for Kinga
This project's purpose is to help Kinga track and plan her diet.

## Features
- N/A

## Building
### Prerequisites
- :white_check_mark: [`go@1.23.1`](https://go.dev/dl/)
- :white_check_mark: [`sqlc@v1.27.0`](https://docs.sqlc.dev/en/latest/overview/install.html)
- :white_check_mark: [`sqlite3`](https://command-not-found.com/sqlite3)

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