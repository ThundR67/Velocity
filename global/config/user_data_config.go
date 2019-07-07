package config

//Here are all the config related to User data like validation stuff

//Loading config manager
var userDataConfigManager = getConfigManager("user_data.config")

var (
	//UserDataConfigAccountStatusActive Is Active Status
	UserDataConfigAccountStatusActive = getStringConfig("userData.userExtraData.accountStatus.Active", userDataConfigManager)
	//UserDataConfigAccountStatusUnactivated Is Unactivated Status
	UserDataConfigAccountStatusUnactivated = getStringConfig("userData.userExtraData.accountStatus.Unactivated", userDataConfigManager)
	//UserDataConfigAccountStatusDeleted Is Deleted Status
	UserDataConfigAccountStatusDeleted = getStringConfig("userData.userExtraData.accountStatus.Deleted", userDataConfigManager)

	//UserDataConfigUserIDCharset Is All Character Which Can Be Used In UserID
	UserDataConfigUserIDCharset = getStringConfig("userData.userID.charset", userDataConfigManager)
	//UserDataConfigUserIDLength Is Length Of UserID
	UserDataConfigUserIDLength = getIntConfig("userData.userID.length", userDataConfigManager)
)

//Validation data (This will change when deploying velocity into production)

//genderTypes
const (
	UserDataConfigGenderTypeMale   = "male"
	UserDataConfigGenderTypeFemale = "female"
	UserDataConfigGenderTypeOther  = "other"
)

//UserDataConfigGenderTypes is all gender types
var UserDataConfigGenderTypes = []string{UserDataConfigGenderTypeMale, UserDataConfigGenderTypeFemale, UserDataConfigGenderTypeOther}

//Length Ranges
var (
	//Main Data
	//UsernameLengthRange is username max and min len
	UserDataConfigUsernameLengthRange = [2]int{3, 16}
	//EmailLengthRange is email max and min len
	UserDataConfigEmailLengthRange = [2]int{5, 100}
	//PasswordLengthRange is password max and min len
	UserDataConfigPasswordLengthRange = [2]int{5, 100}

	//Extra Data
	//FirstNameLengthRange is First max and min len
	UserDataConfigFirstNameLengthRange = [2]int{5, 100}
	//LastNameLengthRange is Last max and min len
	UserDataConfigLastNameLengthRange = [2]int{5, 100}
	//MinimumAge is minimum age required to join
	UserDataConfigMinimumAge = 14
)
