package config

import (
	"os"
	"strings"
	"time"

	goup "github.com/ufoscout/go-up"
)

func getConfigFilePath() string {
	dir, _ := os.Getwd()
	split := strings.Split(dir, "\\")
	if split[len(split)-1] == "user-data-manager" {
		return "../config/main.config"
	}
	return "config/main.config"
}

//Loading Up Config File
var configurations, _ = goup.NewGoUp().
	AddFile(getConfigFilePath(), false).
	Build()

//UserDoesNotExistError Occurs When A User Does Not Exist
type UserDoesNotExistError struct{}

func (userDoesNotExistError UserDoesNotExistError) Error() string {
	return "User Does Not Exist"
}

var (
	//ConfigLogFile Is The Log File
	ConfigLogFile = configurations.GetString("logfile")

	//ConfigMongoDBAddress MongoDb Address
	ConfigMongoDBAddress = configurations.GetString("mongodb.address")
	//ConfigTimeout Is The Timeout While Connecting To MongoDB
	ConfigTimeout = time.Duration(configurations.GetInt("mongodb.timeoutSeconds")) * time.Second
	//ConfigZeroTechhDB Is The Main DB Address
	ConfigZeroTechhDB = configurations.GetString("mongodb.zerotechhDB")
	//ConfigUserDataCollection Is UserData Collection
	ConfigUserDataCollection = configurations.GetString("mongodb.zerotechhDB.UserDataCollection")

	//ConfigUserIDField Is UserID Field
	ConfigUserIDField = configurations.GetString("mongodb.zerotechhDB.UserDataCollection.UserIDField")
	//ConfigUsernameField Is Username Field
	ConfigUsernameField = configurations.GetString("mongodb.zerotechhDB.UserDataCollection.UsernameField")
	//ConfigEmailField Is Email Field
	ConfigEmailField = configurations.GetString("mongodb.zerotechhDB.UserDataCollection.EmailField")
	//ConfigPasswordField Is Password Field
	ConfigPasswordField = configurations.GetString("mongodb.zerotechhDB.UserDataCollection.PasswordField")

	//ConfigUserExtraDataField Is Where Extra Data Is Stored
	ConfigUserExtraDataField = configurations.GetString("mongodb.zerotechhDB.UserDataCollection.UserExtraDataField")

	//ConfigAccountCreationUTCField Is UTC Of When Account Was Created
	ConfigAccountCreationUTCField = configurations.GetString("mongodb.zerotechhDB.UserDataCollection.UserExtraData.AccountCreationUTCField")
	//ConfigAccountStatusField Is Current Status Of Account Field
	ConfigAccountStatusField = configurations.GetString("mongodb.zerotechhDB.UserDataCollection.UserExtraData.AccountStatusField")

	//ConfigAccountStatusActive Is Active Status
	ConfigAccountStatusActive = configurations.GetString("accountStatus.Active")

	//ConfigAccountStatusDeleted Is Deleted Status
	ConfigAccountStatusDeleted = configurations.GetString("accountStatus.Deleted")

	//ConfigUserIDCharset Is All Character Which Can Be Used In UserID
	ConfigUserIDCharset = configurations.GetString("userID.charset")
	//ConfigUserIDLength Is Length Of UserID
	ConfigUserIDLength = configurations.GetInt("userID.length")
)

const (
	//UsernameExistsMsg Is When Username Is Taken
	UsernameExistsMsg = "Username Exists"
	//EmailExistsMsg Is When Email Is Taken
	EmailExistsMsg = "Email Exists"
)
