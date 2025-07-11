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
- Go 1.23.1+
    - TODO go install the project
    - TODO go install sqlc
    - TODO go install goose
    - TODO install `dlv`
- Postgres 15
    - Connection string
- API client (ex: Postman)


# Usage
## First Run
### Create DB
### Create the `.env` file

1. Create a file named `.env` in the root directory and add it to .gitgnore. **DO NOT COMMIT .env**.

1. In `.env`, add the following variables:

    ```shell
    DB_URL="postgres://<username>:<password>@localhost:5432/chirpy?sslmode=disable"
    PLATFORM="dev"
    JWT_SECRET="TODO: Steps to generate"
    POLKA_API_KEY="<random, alphanumeric 32-char fake api key>"
    ```

1. From the root directory run
    
        $ go run .
    
1. `reset` endpoint
`POST http://localhost:8080/api/reset` Deletes **all* users and posts, clearing the database.

1. (Optional) External API client
[Postman Collection](/docs/chirpy.postman_collection.json)

# Endpoints


# Development Notes
## Database Schema Changes
1. Create new [Goose migration](https://github.com/pressly/goose) file, place in [`sql/schema`](sql/schema/).

2. Run migration:

        $ cd sql/schema
        sql/schema $ goose postgres <postgress connection string> up

        # Rollback last migration (if needed)
        sql/schema $ goose postgres <postgress connection string> down

## Database Queries

Use SQLC to generate GO functions from SQL queries:

1. Add plain SQL queries here: [`sql/schema`](sql/queries/)
2. Generate Go SQL functions:

    $ cd <base directory>
    $ sqlc generate

3. Access queries from the `apiConfig` struct:

    ```golang
    func (cfg *apiConfig) myHandler() http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {

            result, err := cfg.db.MySQLFunction(...)

        }
    }
    ```
See the various handler functions for usage examples.

## Debugging
1. Start the [Delve debugger](https://github.com/go-delve/delve):

        $ dlv debug . --headless --listen=:12345 --continue --accept-multiclient

2. Connect the remote server
connect with remote debugger, ex: launch.json (VSCode)
    - Link to launch.json
- Place breakpoint
- Send requests, ex: 

# Tests
- Run unit tests: `$ go test ./...`
    - Unit test helper functions
    - Go table structure testing pattern
        - Slice of `structs` to hold different "cases", followed by actual unit test code
            - Improves test readability and easy to add new cases

# Potential enhancements
- Local Front-end interface
- Fuzz testing with Go's testing libraries
- Implement code test coverage
- Deployment to cloud service; publicly accessible
