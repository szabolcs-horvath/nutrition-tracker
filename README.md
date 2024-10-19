# Nutrition Tracker for Kinga
This project's purpose is to help Kinga track and plan her diet.

## Features:
- N/A

### Prerequisites:
- [ ] `go@1.23.1`
- [ ] `sqlc@v1.27.0`
- [ ] `sqlite3`

### Building
1. #### Create the database
    ```shell
    sqlite3 sqlite/nutrition-tracker.db
    ```
    You can exit from the interactive sqlite console with `.quit`

2. #### Run the migrations
    ```shell
    make migrate-up
    ```
3. #### Build the project
    ```shell
    make build
    ```
