package config

//Here are all the config related to auth server

//Loading config manager
var authServerConfigManager = getConfigManager("auth_server.config")

var (
	//AuthServerConfigIPAddress the ip where auth server will run
	AuthServerConfigIPAddress = getStringConfig("authServer.ipAddress", authServerConfigManager)
	//AuthServerConfigPort is the port where auth server will run
	AuthServerConfigPort = getIntConfig("authServer.port", authServerConfigManager)
	//AuthServerConfigUsernameField Is username field in http request
	AuthServerConfigUsernameField = getStringConfig("authServer.request.usernameField", authServerConfigManager)
	//AuthServerConfigPasswordField Is password field in http request
	AuthServerConfigPasswordField = getStringConfig("authServer.request.passwordField", authServerConfigManager)
	//AuthServerConfigScopesField Is scopes field in request
	AuthServerConfigScopesField = getStringConfig("authServer.request.scopesField", authServerConfigManager)
	//AuthServerConfigErrField is err field of response
	AuthServerConfigErrField = getStringConfig("authServer.response.errField", authServerConfigManager)
	//AuthServerConfigAccessTokenField is field of access token
	AuthServerConfigAccessTokenField = getStringConfig("authServer.response.accessTokenField", authServerConfigManager)
	//AuthServerConfigRefreshTokenField is field of refresh token
	AuthServerConfigRefreshTokenField = getStringConfig("authServer.response.refreshTokenField", authServerConfigManager)
)

const (
	//InternalServerError Is when internal server error
	InternalServerError = "Internal Server Error"
)
