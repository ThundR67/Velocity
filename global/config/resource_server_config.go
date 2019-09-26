package config

var resourceServerConfigManager = getConfigManager("resource_server.config")

var (
	//ResourceSrvConfigJWTHeader is the key to jwt token header
	ResourceSrvConfigJWTHeader = resourceServerConfigManager.GetString("resourceServer.jwtHeader")
)

const (
	//InvalidInputMsg is used when inputs are invalid
	InvalidInputMsg = "Invalid Inputs"
)
