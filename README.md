# Chirpy

A simple REST server written in Golang to practice the basics of developing an API service from scratch.
The server simulates a locally running social media service where "users" can post "chirps".

Guided project using backend developer training site [boot.dev](https://www.boot.dev/lessons/50f37da8-72c0-4860-a7d1-17e4bda5c243).

Concepts covered:
- HTTP `GET`, `POST`, `PUT`, `DELETE` operations
- PostgresSQL acceess and storage
- User authentication
- Endpoint Authorization
- Query parameters
- Middleware functions
- 3rd-party Webhooks
- HTTP unit testing

# Setup
- Go 1.23.1+
- Postgres 15
- API client (ex: Postman)
- .env secrets

# Usage
- go run .
- API client
- `reset` endpoint

# Endpoints


# Development
1. Write db query, if needed
1. `sqlc generate` Go function
1. Create handler
1. Write logic: extract tokens, access db, generate responses, etc.

## Debugging
- `dlv` remote server
- Connect as remote debugger
    - launch.json (VSCode)
- Sen request using api client, ex: Postman
    - Include Postman collection files
    - Scripts to auto save values

# Tests
- go test ./...
- helper functions
- GO table structure testing pattern

# Potential enhancements
- Front-end interface
- Fuzz testing
- Code test coverage
- Deployment to cloud service
