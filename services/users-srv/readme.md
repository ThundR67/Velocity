# Users Service
### This service will handle everything related to user data. This service will add, get and authenticate user data.


## Data
#### This service splits and stores data into three sections.

### Main Data
##### Containing all the essential data of user such as username and email.

### Extra Data
##### Containing non important data such as first and last name.

### Meta Data
##### Containing data which is important to the server such as account status and account creation time.

## User Type
### This is the user type shown in JSON format, this user type is defined in global/config/types.go

```javascript
{
    "_id" : "Users ID Which Is UUID v4",
    //MainData
    "Username" : "A Unique Lowercased Username",
    "Email" : "Unique Email Of The User",
    "Password" : "A hashed password of the user",
    //ExtraData
    "FirstName" : "First name of the user",
    "LastName" : "Last name of the user",]
    "Gender" : "Gender of the user",
    "BirthdayUTC" : "Users birthday timestamp",
    //MetaData
    "AccountCreationUTC" : "Users account creation timestamp",
    "AccountStatus" : "Account status of user",

}
```

## Service Functions
### These are all the functions of this service, defined in proto/users-srv.proto

### Add
#### Add function takes in user data, validates it, adds it to db, and then returnes user id.

### Get
#### Get function takes in user id and returnes user's main data

### GetExtra
#### GetExtra function takes in user id and returns user's extra data

### GetByUsernameOrEmail
#### Returned users main data based on either username or email

### Auth
#### Takes in username and password, then authenticates and returnes user id

### Update
#### Update takes in user id and the update. Then updates user's main data

### UpdateExtra
#### UpdateExtra takes in user id and the update. Then updates user's extra data

### Delete
#### Delete takes user id and marks that user as deleted.