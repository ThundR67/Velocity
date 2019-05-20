package userdatamanager

import (
	"strings"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
)

func timeFromRequest(ts int64) time.Time {
	return time.Unix(ts, 0)
}

func lenCheck(value int, lenRange [2]int) bool {
	return lenRange[0] <= value && value <= lenRange[1]
}

func genderCheck(value string) bool {
	for _, gender := range config.UserDataConfigGenderTypes {
		if value == gender {
			return true
		}
	}
	return false
}

func validateUserData(userData map[string]interface{}) (valid bool) {

	defer func() {
		if r := recover(); r != nil {
			valid = false
		}
	}()

	//checking if username is lower case
	username := userData[config.DBConfigUsernameField].(string)
	if username != strings.ToLower(username) {
		return false
	}

	//Len Checks
	firstNameLen := len(userData[config.DBConfigUserExtraDataField].(map[string]interface{})[config.DBConfigFirstNameField].(string))
	lastNameLen := len(userData[config.DBConfigUserExtraDataField].(map[string]interface{})[config.DBConfigLastNameField].(string))
	gender := userData[config.DBConfigUserExtraDataField].(map[string]interface{})[config.DBConfigGenderField].(string)

	if !lenCheck(len(username), config.UserDataConfigUsernameLengthRange) {
		return false
	} else if !lenCheck(len(userData[config.DBConfigPasswordField].(string)), config.UserDataConfigPasswordLengthRange) {
		return false
	} else if !lenCheck(len(userData[config.DBConfigEmailField].(string)), config.UserDataConfigEmailLengthRange) {
		return false
	} else if !lenCheck(len(userData[config.DBConfigUsernameField].(string)), config.UserDataConfigUsernameLengthRange) {
		return false
	} else if !lenCheck(len(userData[config.DBConfigUsernameField].(string)), config.UserDataConfigUsernameLengthRange) {
		return false
	} else if !lenCheck(firstNameLen, config.UserDataConfigFirstNameLengthRange) {
		return false
	} else if !lenCheck(lastNameLen, config.UserDataConfigLastNameLengthRange) {
		return false
	}

	//gender check
	if !genderCheck(gender) {
		return false
	}

	birthday := timeFromRequest(userData[config.DBConfigUserExtraDataField].(map[string]interface{})[config.DBConfigBirthdayUTCField].(int64))
	if time.Now().Year()-birthday.Year() < config.UserDataConfigMinimumAge {
		return false
	}
	return true
}
