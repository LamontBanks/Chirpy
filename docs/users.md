# API Documentation

## Authentication
- **Type**: Bearer Token / API Key
- **Header**: `Authorization: Bearer <token>` / `Authorization: ApiKey <token>`

## Base URL

http://localhost:8080

## Users

### POST `/api/users/`

Register a user.

**Headers**:
- `Content-Type`: `application/json`

**Parameters**:
- `email` (required): Email Address (string)
- `password` (required): No requirements, but must be non-empty  (string)

**Sample Request**:
```json
POST /api/user
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "my_secure_password_123"
}
```

**Responses**:

- **200 OK**

    User is registered.

```json
{
    "id": "769f7956-2fe7-4f33-989c-1d5f9cc0cf24",
    "email": "user@example.com",
    "created_at": "2025-07-09T10:04:02.824303Z",
    "updated_at": "2025-07-09T10:04:02.824303Z",
    "is_chirpy_red": false
}
```

- `id`: user id (string)
- `email`: user's email (string)
- `created_at`: user registration timestamp (string)
- `updated_at`: timestamp of last user update (example, changing email, password) (string)
- `is_chirpy_red`: `true` is user is a premium Chirpy member, (default: `false`)


- **500 Internal Server Error**

    Attempted to register user with that already exists:
```json
{
    "error": "Unable to create user with email user@example.com"
}
```

Missing required field(s):
```json
{
    "error": "Something went wrong"
}
```

### GET `/api/users/`

Retrieve all users

**Parameters**:
- `id` (path, required): User ID (uuid)

**Headers**:
- `Content-Type`: application/json

**Sample Request**
```
GET http://localhost:8080/api/users
```

**Responses**:

- **200 OK**
```json
[
    {
        "id": "769f7956-2fe7-4f33-989c-1d5f9cc0cf24",
        "email": "user@example.com",
        "created_at": "2025-07-09T10:04:02.824303Z",
        "updated_at": "2025-07-09T10:04:02.824303Z",
        "is_chirpy_red": false
    },
    {
        "id": "1937f0c5-76eb-490a-85a5-8bdde376e4a4",
        "email": "test@example.com",
        "created_at": "2025-07-09T10:46:32.460291Z",
        "updated_at": "2025-07-09T10:46:32.460291Z",
        "is_chirpy_red": false
    }
]
```


### GET /users/{id}
Retrieve a specific user by ID

**Parameters**:
- `id` (path, required): User ID (UUID)

**Headers**:
- `Authorization`: Bearer token (required)
- `Content-Type`: application/json

**Responses**:
- **200 OK**
```json
{
"id": 123,
"name": "John Doe",
"email": "john@example.com"
}
```

