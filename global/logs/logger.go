package logger

import (
	"github.com/SonicRoshan/Velocity/global/config"
	log "github.com/jex-lin/golang-logger"
)

//getLogFilePath returnes relative path of log file
func getLogFilePath(fileName string) string {
	return "D:/VelocityLogs/" + fileName
}

//GetLogger returns a logger
func GetLogger(fileName string) *log.Log {
	output := log.NewLogFile(getLogFilePath(fileName))
	output.SetLevel(config.LogConfigLogLevel)
	return output
}
