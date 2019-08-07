package config

import (
	"fmt"
	"os"
	"strings"

	goup "github.com/ufoscout/go-up"
)

//getConfigFilePath returnes relative path of config file
func getConfigFilePath(fileName string) string {
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
	filepath += fmt.Sprintf("global\\config\\config\\%s", fileName)
	return filepath
}

//getConfigManager returnes a config manager based fileName
func getConfigManager(fileName string) goup.GoUp {
	configManager, err := goup.NewGoUp().
		AddFile(getConfigFilePath(fileName), false).
		Build()

	if err != nil {
		panic(fmt.Sprintf("Loading Config File %s Returned Error %s", fileName, err.Error()))
	} else if configManager == nil {
		panic(fmt.Sprintf("Config File %s Is Nil", fileName))
	}
	return configManager
}

//getStringConfig gets a configuration with string value
func getStringConfig(name string, configManager goup.GoUp) string {
	value, err := configManager.GetStringOrFail(name)
	if err != nil {
		panic(fmt.Sprintf("Cannot Load %s Config From Config File, With Error %s", name, err.Error()))
	}
	return value
}

//getStringSliceConfig gets a configuration with string slice value
func getStringSliceConfig(name string, configManager goup.GoUp) []string {
	value, err := configManager.GetStringSliceOrFail(name, ",")
	if err != nil {
		panic(fmt.Sprintf("Cannot Load %s Config From Config File, With Error %s", name, err.Error()))
	}
	return value
}

//getIntConfig gets a configuration with int value
func getIntConfig(name string, configManager goup.GoUp) int {
	//TODO Add logging
	value, err := configManager.GetIntOrFail(name)
	if err != nil {
		panic(fmt.Sprintf("Cannot Load %s Config From Config File, With Error %s", name, err.Error()))
	}
	return value
}

//getBoolConfig gets a configuration with bool value
func getBoolConfig(name string, configManager goup.GoUp) bool {
	//TODO Add logging
	value, err := configManager.GetBoolOrFail(name)
	if err != nil {
		panic(fmt.Sprintf("Cannot Load %s Config From Config File, With Error %s", name, err.Error()))
	}
	return value
}

//CustomError is custom error type
type CustomError struct {
	msg string
}

func (customError CustomError) Error() string {
	return customError.msg
}

func generateErrorWithMsg(msg string) error {
	customError := CustomError{msg: msg}
	return customError
}
