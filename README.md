# Chirpy

A simple REST server written from scratch in Golang to practice the basics of developing an API service.
The server simulates a locally running "BlueSky/Twitter-like" API where users, register, login, post, and view messages.

Guided project using backend developer training site [boot.dev](https://www.boot.dev/lessons/50f37da8-72c0-4860-a7d1-17e4bda5c243).

Concepts covered:
- `GET`, `POST`, `PUT`, `DELETE` operations
- PostgresSQL acceess and storage
- Database migrations
- User authentication
- Endpoint authorization
- Query parameters
- Middleware functions
- 3rd-party webhooks
- HTTP unit testing

# Setup
- Install Go 1.23.1+:

- Install project libraries



- Install [SQLC](https://github.com/sqlc-dev/sqlc):

    `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

- Install [Goose](https://github.com/pressly/goose):

    `go install github.com/pressly/goose/v3/cmd/goose@latest`

- Install the `Delve` debugger: https://github.com/go-delve/delve/tree/master/Documentation/installation

- (Optional) Install an API client (ex: Postman)
    - [Postman collection](/docs/chirpy.postman_collection.json)

# Usage

## Preliminary
### Create DB
1. Install PostgresSQL (Mac): 

    `brew install postgresql@15`

1. Run in the background:

    `brew services start postgresql@15`

1. Use a client to connect (ex: `psql`):

    `psql postgres`

1. Create the blank database:

    ```sql
    CREATE DATABASE chirpy;
    ```

### Create the `.env` file

1. Create a file named `.env` in the root directory and add it to .gitgnore. **DO NOT COMMIT .env**.

1. In `.env`, add the following variables:

    ```shell
    DB_URL="postgres://<username>:<password>@localhost:5432/chirpy?sslmode=disable"
    PLATFORM="dev"
    JWT_SECRET="TODO: Steps to generate"
    POLKA_API_KEY="<random, alphanumeric 32-char fake api key>"
    ```

## Start server
1. From the root directory: 

        go run .
    
1. Run `POST http://localhost:8080/api/reset` to delete **all* users and posts, clearing the database
    - Can also use during development testing

1. Send requests to `http://localhost:8080/<endpoint>`

# Endpoints

See full [documentation](/docs/).

# Development Notes
## Database Schema Changes
1. Create new [Goose migration](https://github.com/pressly/goose) file, place in [`sql/schema`](sql/schema/).

2. Run migration:

    ```shell
    cd sql/schema
    goose postgres <postgress connection string> up

    # Rollback last migration (if needed)
    goose postgres <postgress connection string> down
    ```

## Database Queries
Use SQLC to generate GO functions from SQL queries:

1. Add plain SQL queries here: [`sql/schema`](sql/queries/)
2. Generate Go functions:

    ```shell
    cd <base directory>
    sqlc generate
    ````

3. Access queries from the `apiConfig` struct:

    ```golang
    func (cfg *apiConfig) myHandler() http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {

            result, err := cfg.db.MySQLFunction(...)

        }
    }
    ```

See the existing handler functions for more examples.

# Handlers

Handlers are set in `main.go` with implementation in seperate files.

## Debugging
1. Start the [Delve debugger](https://github.com/go-delve/delve):

    ```shell
    dlv debug . --headless --listen=:12345 --continue --accept-multiclient
    ```

2. Connect to the debugger
    - Working VSCode `launch.json` config:

    ```json
    {
        "version": "0.2.0",
        "configurations": [
            {
                "name": "Connect to external session",
                "type": "go",
                "debugAdapter": "dlv-dap",
                "request": "attach",
                "mode": "remote",
                "port": 12345
                // "host": "127.0.0.1", // can skip for localhost
            }
        ]
    }
    ```

3. Set breakpoints in code
4. Send requests to server (ex: using Postman)
5. Step debug

# Tests
- Unit tests: `go test ./...`

# Potential enhancements
- Implement code test coverage
- Experiment with fuzz testing with Go's testing libraries
- Local Front-end interface
- Deployment to cloud service; publicly accessible
