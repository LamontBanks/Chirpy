# Webhooks

"Chirps" are Tweet-like messages user can post, view, and delete.

Concepts demonstrated:
- Create, Read, and Delete data
- Filter and sort data based on URL query parameters

### `GET http://localhost:8080/api/chirps`

Get all chirps, sorted from oldest to newest created (default)

### Responses

Sample Response

```json
[
    {
        "id": "d531eaa4-e7dd-488a-8a9f-18750fd17969",
        "created_at": "2025-07-11T18:15:01.949959Z",
        "updated_at": "2025-07-11T18:15:01.949959Z",
        "body": "first chirp",
        "user_id": "04af8611-66f2-4ab5-bd47-e2255a398b7e"
    },
    {
        "id": "dd2af79b-cf18-4a67-93c5-c5ef58ac9a5b",
        "created_at": "2025-07-11T18:15:04.771555Z",
        "updated_at": "2025-07-11T18:15:04.771555Z",
        "body": "second chirp",
        "user_id": "04af8611-66f2-4ab5-bd47-e2255a398b7e"
    },
    {
        "id": "697a63a5-1af1-4fe1-aa01-1ca154cc6663",
        "created_at": "2025-07-11T18:16:10.806804Z",
        "updated_at": "2025-07-11T18:16:10.806804Z",
        "body": "hey everyone",
        "user_id": "04af8611-66f2-4ab5-bd47-e2255a398b7e"
    },
    {
        "id": "8fcc55b0-3c81-4bbc-a6f3-b38a15ff26c8",
        "created_at": "2025-07-11T18:16:20.67289Z",
        "updated_at": "2025-07-11T18:16:20.672891Z",
        "body": "going on a trip this weekend",
        "user_id": "04af8611-66f2-4ab5-bd47-e2255a398b7e"
    }
]
```

- `HTTP 204`
    - Payment info processed and user upgraded to Chirpy Red

    Also
    cccccbrfulknrbvdenctlegihkndfejglbbfjflbdrfb
    
    - Not a `"user.upgraded"` event

- `HTTP 400 Bad Request`
    - Missing request body elements

- `HTTP 401 Unauthorized`
    - Error processing Polka API Key

- `HTTP 404`
    - User not found or failed to update user
