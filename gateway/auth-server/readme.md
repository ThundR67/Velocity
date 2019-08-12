# Auth Server

Auth server will authorize the client and provide them with JWT


## Routes

### /sign-in
This will take username and password, then return access and refresh tokens.

### /sign-in-fresh
This will take username and password, then return fresh access token.

### /sign-up
This takes user data, adds that user, then returns access and refresh tokens.

### /refresh
This takes refresh token, validates it, then returns new access and refresh token.

### /verify-email
Takes in verification code and verfies the email.