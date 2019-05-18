package jwtmanager

//InvalidTokenError is when when a token is invalid
type InvalidTokenError struct{}

func (ite InvalidTokenError) Error() string { return "Invalid Token" }

//TokenExpiredError is when when a token expired
type TokenExpiredError struct{}

func (tee TokenExpiredError) Error() string { return "Token Has Expired" }

//InvalidScopesError is when scopes are invalid or not allowed
type InvalidScopesError struct{}

func (ise InvalidScopesError) Error() string { return "Invalid Scopes" }
