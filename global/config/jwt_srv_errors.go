package config

var (
	//InvalidTokenError is when token is invalid
	InvalidTokenError = generateErrorWithMsg("Invalid Token")
	//TokenExpiredError is when token has expired
	TokenExpiredError = generateErrorWithMsg("Token Has Expired")
	//InvalidScopesError is when scopes are invalid
	InvalidScopesError = generateErrorWithMsg("Invalid Scopes")
)
