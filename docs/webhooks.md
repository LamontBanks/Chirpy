# Webhooks

"Polka" is an imaginary 3rd-party payment service, used to demonstrate basic webhook functionality.

Concepts demonstrated:
- Processing a vendor supplied payload: `user.upgraded` event and additional `data` block.

## `POST http://localhost:8080/api/polka/webhooks`

Polka triggers an update to a user's Chirpy Red premium status after user completes a "payment".

### Headers
- `Authorization`: `ApiKey <Polka API Key>`

### Request Body
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "string"
  }
}
```

### Responses

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
