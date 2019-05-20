package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/SonicRoshan/Velocity/global/config"
	log "github.com/jex-lin/golang-logger"
)

//getLogFilePath returnes relative path of log file
func getLogFilePath(fileName string) string {
	workingDir, _ := os.Getwd()
	workingDirSplit := strings.Split(workingDir, "\\")
	filepath := ""
	velocityCame := false
	for _, path := range workingDirSplit {
		velocityCame = path == "Velocity"
		if velocityCame {
			filepath += path + "\\"
			break
		}
		filepath += path + "\\"
	}
	filepath += fmt.Sprintf("global\\logs\\logs\\%s", fileName)
	return filepath
}

//GetLogger returns a logger
func GetLogger(fileName string) *log.Log {
	output := log.NewLogFile(getLogFilePath(fileName))
	output.SetLevel(config.LogConfigLogLevel)
	return output
}
