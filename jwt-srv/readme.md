# JWT Service
#### This service will handle every tasks related to Json Web Token.
#### This will create, authenticate, etc JWT.
#### This will also be managing scopes.


## Schemes

### Scheme Of Acces Token Claim
```python
{
    "userIdentity" : {User ID},
    "fresh" : false,
    "scopes" : [Array Of Scopes],
    "creationUTC" : {UTC Of When This JWT Was Created},
    "expirationUTC" : {UTC Of When This JWT Will Expire},
}
```

### Scheme Of Fresh Access Token
##### A fresh acces token is extremely short lived access token. This can be used for example when the user wants to change their password. To confirm and prove their identity we can ask them to give their username and password again. We then issue them a fresh access token instead of regular access token, and only with the fresh acces token, they will be allowed to change their password. And fresh access token can't be refreshed with refresh token.
```python
{
    "userIdentity" : {User ID},
    "fresh" : true,
    "creationUTC" : {UTC Of When This JWT Was Created},
    "expirationUTC" : {UTC Of When This JWT Will Expire},
}
```

### Scheme Of Refresh Token Claim
```python
{
    "userIdentity" : {User ID},
    "scopes" : [Array Of Scopes],
    "creationUTC" : {UTC Of When This JWT Was Created},
    "expirationUTC" : {UTC Of When This JWT Will Expire},
}
```

