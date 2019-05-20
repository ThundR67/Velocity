package config

import "time"

//Here are all the config related to Database

//Loading config manager
var dbConfigManager = getConfigManager("database.config")

var (

	//DBConfigMongoDBAddress MongoDb Address
	DBConfigMongoDBAddress = getStringConfig("database.address", dbConfigManager)
	//DBConfigTimeout Is The Timeout While Connecting To MongoDB
	DBConfigTimeout = time.Duration(getIntConfig("database.timeoutSeconds", dbConfigManager)) * time.Second
	//DBConfigZeroTechhDB Is The ZeroTechhDB Name
	DBConfigZeroTechhDB = getStringConfig("database.zerotechhDB", dbConfigManager)
	//DBConfigUserDataCollection Is UserData Collection
	DBConfigUserDataCollection = getStringConfig("database.zerotechhDB.UserDataCollection", dbConfigManager)

	//DBConfigUserIDField Is UserID Field
	DBConfigUserIDField = getStringConfig("database.zerotechhDB.UserDataCollection.UserIDField", dbConfigManager)
	//DBConfigUsernameField Is Username Field
	DBConfigUsernameField = getStringConfig("database.zerotechhDB.UserDataCollection.UsernameField", dbConfigManager)
	//DBConfigEmailField Is Email Field
	DBConfigEmailField = getStringConfig("database.zerotechhDB.UserDataCollection.EmailField", dbConfigManager)
	//DBConfigPasswordField Is Password Field
	DBConfigPasswordField = getStringConfig("database.zerotechhDB.UserDataCollection.PasswordField", dbConfigManager)

	//DBConfigUserExtraDataField Is Where Extra Data Is Stored
	DBConfigUserExtraDataField = getStringConfig("database.zerotechhDB.UserDataCollection.UserExtraDataField", dbConfigManager)

	//DBConfigBirthdayUTCField is the birthday field
	DBConfigBirthdayUTCField = getStringConfig("database.zerotechhDB.UserDataCollection.UserExtraData.BirthdayUTCField", dbConfigManager)
	//DBConfigGenderField is the gender field
	DBConfigGenderField = getStringConfig("database.zerotechhDB.UserDataCollection.UserExtraData.GenderField", dbConfigManager)
	//DBConfigFirstNameField is first name field
	DBConfigFirstNameField = getStringConfig("database.zerotechhDB.UserDataCollection.UserExtraData.FirstNameField", dbConfigManager)
	//DBConfigLastNameField is last name field
	DBConfigLastNameField = getStringConfig("database.zerotechhDB.UserDataCollection.UserExtraData.LastNameField", dbConfigManager)
	//DBConfigAccountCreationUTCField Is UTC Of When Account Was Created
	DBConfigAccountCreationUTCField = getStringConfig("database.zerotechhDB.UserDataCollection.UserExtraData.AccountCreationUTCField", dbConfigManager)
	//DBConfigAccountStatusField Is Current Status Of Account Field
	DBConfigAccountStatusField = getStringConfig("database.zerotechhDB.UserDataCollection.UserExtraData.AccountStatusField", dbConfigManager)
)
