package userdatamanager

import (
	"strconv"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/asaskevich/govalidator"
)

func isGender(genderInput string) bool {
	for _, gender := range config.UserDataConfigGenderTypes {
		if gender == genderInput {
			return true
		}
	}
	return false
}

func handlePanic() bool {
	r := recover()
	return r == nil
}

func timeFromRequest(ts string) time.Time {
	i, _ := strconv.ParseInt(ts, 10, 64)
	return time.Unix(i, 0)
}

func validateUserMainData(userMainData map[string]interface{}) (valid bool) {
	defer func() { valid = handlePanic() }()
	username := userMainData[config.DBConfigUsernameField].(string)
	password := userMainData[config.DBConfigPasswordField].(string)
	email := userMainData[config.DBConfigEmailField].(string)
	if govalidator.HasUpperCase(username) {
		return false
	} else if !govalidator.InRange(len(username), config.UserDataConfigUsernameLengthRange[0], config.UserDataConfigUsernameLengthRange[1]) {
		return false
	} else if !govalidator.IsExistingEmail(email) {
		return false
	} else if !govalidator.InRange(len(password), config.UserDataConfigPasswordLengthRange[0], config.UserDataConfigPasswordLengthRange[1]) {
		return false
	}
	return true
}

func validateUserExtraData(userExtraData map[string]interface{}) (valid bool) {
	defer func() { valid = handlePanic() }()
	firstname := userExtraData[config.DBConfigFirstNameField].(string)
	lastname := userExtraData[config.DBConfigLastNameField].(string)
	gender := userExtraData[config.DBConfigGenderField].(string)
	birthday := timeFromRequest(userExtraData[config.DBConfigBirthdayUTCField].(string))
	age := time.Now().Year() - birthday.Year()

	if !govalidator.InRange(len(firstname), config.UserDataConfigFirstNameLengthRange[0], config.UserDataConfigFirstNameLengthRange[1]) {
		return false
	} else if !govalidator.InRange(len(lastname), config.UserDataConfigLastNameLengthRange[0], config.UserDataConfigLastNameLengthRange[1]) {
		return false
	} else if !isGender(gender) {
		return false
	} else if age < config.UserDataConfigMinimumAge {
		return false
	}
	return true
}
