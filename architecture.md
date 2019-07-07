# Velocity Architecture
### Architecture of Velocity
## Gateway
---
#### This is the main part of velocity
#### This is an Oauth2 server which will talk with the client, Only gateway will talk to clients.
#### The gateway will also provide resources

## user-data-srv
---
#### User data service will handle tasks like adding user, getting user data, authenticating user, etc.
#### Any task related to handling user data will go to user data service.

