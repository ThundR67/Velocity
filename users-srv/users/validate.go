package users

import (
	"strings"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/asaskevich/govalidator"
)

var validateLog = logger.GetLogger("users_data_validater.log")

func isGender(genderInput string) bool {
	for _, gender := range config.UserDataConfigGenderTypes {
		if gender == genderInput {
			return true
		}
	}
	return false
}

func isValid(main, extra config.UserType) bool {
	return validateUserMainData(main) && validateUserExtraData(extra)
}

func panicOccured() bool {
	r := recover()
	return r != nil
}

func validateUserMainData(userMainData config.UserType) (valid bool) {
	defer func() {
		if panicOccured() {
			valid = false
		}
	}()

	validateLog.Debugf("Got User Main Data %+v", userMainData)

	if govalidator.HasUpperCase(userMainData.Username) {
		validateLog.Debugf("Username %s Has UpperCase Characters In It", userMainData.Username)
		return false
	} else if !govalidator.InRange(len(userMainData.Username), config.UserDataConfigUsernameLengthRange[0], config.UserDataConfigUsernameLengthRange[1]) {
		validateLog.Debug("Username Has Invalid Length")
		return false
	} else if !govalidator.IsEmail(userMainData.Email) {
		validateLog.Debug("Email Is Invalid")
		return false
	} else if !govalidator.InRange(len(userMainData.Password), config.UserDataConfigPasswordLengthRange[0], config.UserDataConfigPasswordLengthRange[1]) {
		validateLog.Debug("Password Has Invalid Length")
		return false
	} else if strings.Contains(userMainData.Username, " ") {
		validateLog.Debug("Username Has Spaces")
		return false
	} else if strings.Contains(userMainData.Email, " ") {
		validateLog.Debug("Email Has Spaces")
		return false
	}

	return true
}

func validateUserExtraData(userExtraData config.UserType) (valid bool) {
	defer func() {
		if panicOccured() {
			valid = false
		}
	}()

	validateLog.Debugf("Got User Extra Data %+v", userExtraData)

	birthday := time.Unix(userExtraData.BirthdayUTC, 0)
	age := time.Now().Year() - birthday.Year()

	if !govalidator.InRange(len(userExtraData.FirstName), config.UserDataConfigFirstNameLengthRange[0], config.UserDataConfigFirstNameLengthRange[1]) {
		validateLog.Debug("Firstname Length Is Invalid")
		return false
	} else if !govalidator.InRange(len(userExtraData.LastName), config.UserDataConfigLastNameLengthRange[0], config.UserDataConfigLastNameLengthRange[1]) {
		validateLog.Debug("Lastname Length Is Invalid")
		return false
	} else if !isGender(userExtraData.Gender) {
		validateLog.Debug("Gender Is Invalid")
		return false
	} else if age < config.UserDataConfigMinimumAge {
		validateLog.Debug("Age Is Invalid")
		return false
	}
	return true
}
