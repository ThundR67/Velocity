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
	return validateUserMainData(main, true) && validateUserExtraData(extra, true)
}

func panicOccured() bool {
	r := recover()
	return r != nil
}

func validateUserMainData(userMainData config.UserMain, strict bool) (valid bool) {
	defer func() {
		if panicOccured() && strict {
			valid = false
		}
	}()

	validateLog.Debug("Got User Main Data", zap.Any("Main Data", userMainData))

	if govalidator.HasUpperCase(userMainData.Username) {
		validateLog.Debug(
			"Username Has UpperCase Characters In It",
			zap.String("Username", userMainData.Username),
		)
		valid = false
		return
	} else if !govalidator.InRange(len(userMainData.Username), config.UserDataConfigUsernameLengthRange[0], config.UserDataConfigUsernameLengthRange[1]) {
		validateLog.Debug("Username Has Invalid Length")
		valid = false
		return
	} else if !govalidator.IsEmail(userMainData.Email) {
		validateLog.Debug("Email Is Invalid")
		valid = false
		return
	} else if !govalidator.InRange(len(userMainData.Password), config.UserDataConfigPasswordLengthRange[0], config.UserDataConfigPasswordLengthRange[1]) {
		validateLog.Debug("Password Has Invalid Length")
		valid = false
		return
	} else if strings.Contains(userMainData.Username, " ") {
		validateLog.Debug("Username Has Spaces")
		valid = false
		return
	} else if strings.Contains(userMainData.Email, " ") {
		validateLog.Debug("Email Has Spaces")
		valid = false
		return
	}

	valid = true
	return
}

func validateUserExtraData(userExtraData config.UserExtra, strict bool) (valid bool) {
	defer func() {
		if panicOccured() && strict {
			valid = false
		}
	}()

	validateLog.Debug("Got User Extra Data", zap.Any("Extra Data", userExtraData))

	birthday := time.Unix(userExtraData.BirthdayUTC, 0)
	age := time.Now().Year() - birthday.Year()

	if !govalidator.InRange(len(userExtraData.FirstName), config.UserDataConfigFirstNameLengthRange[0], config.UserDataConfigFirstNameLengthRange[1]) {
		validateLog.Debug("Firstname Length Is Invalid")
		valid = false
		return
	} else if !govalidator.InRange(len(userExtraData.LastName), config.UserDataConfigLastNameLengthRange[0], config.UserDataConfigLastNameLengthRange[1]) {
		validateLog.Debug("Lastname Length Is Invalid")
		valid = false
		return
	} else if !isGender(userExtraData.Gender) {
		validateLog.Debug("Gender Is Invalid")
		valid = false
		return
	} else if age < config.UserDataConfigMinimumAge {
		validateLog.Debug("Age Is Invalid")
		valid = false
		return
	}

	valid = true
	return
}
