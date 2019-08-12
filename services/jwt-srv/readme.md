# JWT Service
This service will handle taks related to JSON Web Tokens. This service will create and validate JWTs.

## Token Types

### Access token
This is the token which will allow the client to get data from resource server.

### Fresh Access Token
This is simiral to access token, but extremely short lived. Tasks such as changing the password, would client to have a fresh access token.

### Refresh Token
This token can be used to refresh access token (NOT FRESH ACCESS TOKEN).

## Claims
These are the claims in every token, defined at global/config/types.go

```javascript
{
    "UserIdentity" : "An identity of user"
    "IsFresh" : "Boolean value saying if the token is fresh toke",      
    "IsRefresh" : "Booleean value saying if the token is refresh token",
    "Scopes" : "All the scopes",       
    "CreationUTC"  : "Timestamp of when the token was created",
    "ExpirationUTC" : "Timestamp of when the token will expire,
}
```

## Service Functions
### these are the functions of the service, also defined in proto/jwt-srv.proto

### FreshToken
#### Takes in user identity and returnes fresh access token.

### AccessAndRefreshTokens
#### Takes in user identity and scopes, then returnes access and refresh tokens.

### RefreshTokens
#### Takes in refresh token, validates it, then returned new access and refresh token.

### ValidateToken
#### Takes in token and token type, then validates it.
