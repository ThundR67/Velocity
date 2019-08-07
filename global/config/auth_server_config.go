package config

//Here are all the config related to auth server

//Loading config manager
var authServerConfigManager = getConfigManager("auth_server.config")

var (
	//AuthServerConfigIPAddresses the ips where auth server will run
	AuthServerConfigIPAddresses = getStringSliceConfig("authServer.ipAddresses", authServerConfigManager)
	//AuthServerConfigPort is the port where auth server will run
	AuthServerConfigPort = getIntConfig("authServer.port", authServerConfigManager)
	//AuthServerConfigTokenField Is token field in request
	AuthServerConfigTokenField = getStringConfig("authServer.request.tokenField", authServerConfigManager)
	//AuthServerConfigScopesField Is scopes field in request
	AuthServerConfigScopesField = getStringConfig("authServer.request.scopesField", authServerConfigManager)
	//AuthServerConfigErrField is err field of response
	AuthServerConfigErrField = getStringConfig("authServer.response.errField", authServerConfigManager)
	//AuthServerConfigAccessTokenField is field of access token
	AuthServerConfigAccessTokenField = getStringConfig("authServer.response.accessTokenField", authServerConfigManager)
	//AuthServerConfigRefreshTokenField is field of refresh token
	AuthServerConfigRefreshTokenField = getStringConfig("authServer.response.refreshTokenField", authServerConfigManager)
	//AuthServerConfigShowError should server show error
	AuthServerConfigShowError = getBoolConfig("authServer.response.showError", authServerConfigManager)
)

const (
	//InternalServerError Is when internal server error
	InternalServerError = "Internal Server Error"
)
