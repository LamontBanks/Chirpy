# Chirpy

A simple REST server written from scratch in Golang to practice the basics of developing an API service.
The server simulates a locally running "BlueSky/Twitter-like" API where users, register, login, post, and view messages.

Guided project using backend developer training site [boot.dev](https://www.boot.dev/lessons/50f37da8-72c0-4860-a7d1-17e4bda5c243).

Concepts covered:

- HTTP `GET`, `POST`, `PUT`, `DELETE` operations
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
- Postgres 15
    - Connection string
- API client (ex: Postman)
- .env secrets

# Usage
1. go run .
2. API client
    - Useful Postman script
3. `reset` endpoint

# Endpoints


# Development
1. Write db query, if needed
1. `sqlc generate` Go function
1. Create handler
1. Write logic: extract tokens, access db, generate responses, etc.

## Debugging
- `dlv` remote server, then connect with remote debugger
    - Ex: launch.json (VSCode)
- API client, ex: Postman
    - Include Postman collection files
    - Scripts to auto save values

# Tests
- Run unit tests: `$ go test ./...`
    - Unit test helper functions
    - Go table structure testing pattern
        - Slice of `structs` to hold different "cases", followed by actual unit test code
            - Improves test readability and easy to add new cases

# Potential enhancements
- Front-end interface
- Fuzz testing with Go's testing libraries
- Implement code test coverage
- Deployment to cloud service; publicly accessible
