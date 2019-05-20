package config

//Here are all the config related to logging

//Loading config manager
var logConfigManager = getConfigManager("log.config")

var (
	//LogConfigLogLevel is the logging level
	LogConfigLogLevel = getStringConfig("loggingLevel", logConfigManager)
)
