package users

import (
	"strings"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

func isGender(genderInput string) bool {
	for _, gender := range config.UserDataConfigGenderTypes {
		if gender == genderInput {
			return true
		}
	}
	return false
}

func isValid(main config.UserMain, extra config.UserExtra) bool {
	return validateUserMainData(main) && validateUserExtraData(extra)
}

func validateUserMainData(userMainData config.UserMain) (valid bool) {
	validateLog.Debug("Got User Main Data", zap.Any("Main Data", userMainData))

	/*
		TODO Change the way main and extra data is validated
		by creating a new library which can help do it.
		Something like cerberus for python.
	*/

	if govalidator.HasUpperCase(userMainData.Username) {
		validateLog.Debug("Username Has UpperCase Characters In It")
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
	}

	return true
}

func validateUserExtraData(userExtraData config.UserExtra) (valid bool) {
	validateLog.Debug("Got User Extra Data", zap.Any("Extra Data", userExtraData))

	birthday := time.Unix(userExtraData.BirthdayUTC, 0)
	age := time.Now().Year() - birthday.Year()
	validateLog.Debug("Got Age", zap.Int("Age", age))

	if !govalidator.InRange(len(userExtraData.FirstName), config.UserDataConfigFirstNameLengthRange[0], config.UserDataConfigFirstNameLengthRange[1]) {
		validateLog.Debug("Firstname Length Is Invalid")
		return false
	} else if !govalidator.InRange(len(userExtraData.LastName), config.UserDataConfigLastNameLengthRange[0], config.UserDataConfigLastNameLengthRange[1]) {
		validateLog.Debug("Lastname Length Is Invalid")
		return false
	} else if strings.Contains(userExtraData.FirstName, " ") {
		validateLog.Debug("First Name Contain Spaces")
		return false
	} else if strings.Contains(userExtraData.LastName, " ") {
		validateLog.Debug("Last Name Contain Spaces")
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
