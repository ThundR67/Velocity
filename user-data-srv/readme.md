# User Data Service
### This service handles all the user data inside the database, like adding user, getting user data, authentication user, etc.

---
## Proto
##### The proto file is in the proto file along with compiled version of that proto file in golang
---
## user-data-manager
##### This folder contains the low level user data manager struct which does the actual operations

---
## All The Function

### AddUser
##### This function takes in user data, then it creates a unique 64 char long id for the user. It then adds user to db and returns the userID. It will also return a message if username or email already exists

### GetUser
##### This function takes a userID then returns the user data
---
### GetUserByUsernameOrEmail
##### This function takes either username or email then returns user data
---
### AuthUser
##### This function takes either username or email and a password then returns true if the credentials are true. Currently the passwords are stored raw, soon hashing will be added to it
----

### UpdateUser
##### This function takes userID, a field, and new value. Then it will update that user's field with new value provided
----

### DeleteUser
##### This funtion takes userID, then it will mark that user as deleted