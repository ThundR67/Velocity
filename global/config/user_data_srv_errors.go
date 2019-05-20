package config

var (
	//UserDoesNotExistError is when user does not exists
	UserDoesNotExistError = generateErrorWithMsg("User Does Not Exists")
	//InvalidPasswordError is when password is invalid
	InvalidPasswordError = generateErrorWithMsg("Invalid Password")
	//UsernameExistError is when username is already taken
	UsernameExistError = generateErrorWithMsg("Username Already Exists")
	//EmailExistError is when email is already taken
	EmailExistError = generateErrorWithMsg("Email Already Exists")
	//InvalidUsernameAndEmailError is when username and email is invalid
	InvalidUsernameAndEmailError = generateErrorWithMsg("Invalid Username And Error")
	//InvalidAuthDataError is when auth data is incorrect
	InvalidAuthDataError = generateErrorWithMsg("Invalid Authentication Credentials")
	//InvalidUserDataError is when user data is not valid
	InvalidUserDataError = generateErrorWithMsg("Invalid User Data")
	//InvalidUsernameOrEmailError is when either username or email is invalid
	InvalidUsernameOrEmailError = generateErrorWithMsg("Invalid Username Or Email Error")
)
