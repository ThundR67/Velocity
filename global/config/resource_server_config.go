package config

var resourceServerConfigManager = getConfigManager("resource_server.config")

var (
	//ResourceSrvConfigJWTHeader is the key to jwt token header
	ResourceSrvConfigJWTHeader = resourceServerConfigManager.GetString("resourceServer.jwtHeader")
	//ResourceSrvConfigUserIDKey is the key to user id
	ResourceSrvConfigUserIDKey = resourceServerConfigManager.GetString("resourceServer.userIDKey")
	//ResourceSrvConfigScopesKey is the key to scopes
	ResourceSrvConfigScopesKey = resourceServerConfigManager.GetString("resourceServer.scopesKey")
)

const (
	//InvalidInputMsg is used when inputs are invalid
	InvalidInputMsg = "Invalid Inputs"
)
