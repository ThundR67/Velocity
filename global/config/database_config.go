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
	//DBConfigUserExtraDataCollection Is Where Extra Data Is Stored
	DBConfigUserExtraDataCollection = getStringConfig("database.zerotechhDB.UserExtraDataCollection", dbConfigManager)
	//DBConfigUserMetaDataCollection Is Where User's meta Data Is Stored
	DBConfigUserMetaDataCollection = getStringConfig("database.zerotechhDB.UserMetaDataCollection", dbConfigManager)

	//DBConfigUserIDField Is UserID Field
	DBConfigUserIDField = getStringConfig("database.zerotechhDB.UserDataCollection.UserIDField", dbConfigManager)
	//DBConfigUsernameField Is Username Field
	DBConfigUsernameField = getStringConfig("database.zerotechhDB.UserDataCollection.UsernameField", dbConfigManager)
	//DBConfigEmailField Is Email Field
	DBConfigEmailField = getStringConfig("database.zerotechhDB.UserDataCollection.EmailField", dbConfigManager)
	//DBConfigPasswordField Is Password Field
	DBConfigPasswordField = getStringConfig("database.zerotechhDB.UserDataCollection.PasswordField", dbConfigManager)

	//DBConfigBirthdayUTCField is the birthday field
	DBConfigBirthdayUTCField = getStringConfig("database.zerotechhDB.UserExtraDataCollection.BirthdayUTCField", dbConfigManager)
	//DBConfigGenderField is the gender field
	DBConfigGenderField = getStringConfig("database.zerotechhDB.UserExtraDataCollection.GenderField", dbConfigManager)
	//DBConfigFirstNameField is first name field
	DBConfigFirstNameField = getStringConfig("database.zerotechhDB.UserExtraDataCollection.FirstNameField", dbConfigManager)
	//DBConfigLastNameField is last name field
	DBConfigLastNameField = getStringConfig("database.zerotechhDB.UserExtraDataCollection.LastNameField", dbConfigManager)

	//DBConfigAccountCreationUTCField Is UTC Of When Account Was Created
	DBConfigAccountCreationUTCField = getStringConfig("database.zerotechhDB.UserMetaDataCollection.AccountCreationUTCField", dbConfigManager)
	//DBConfigAccountStatusField Is Current Status Of Account Field
	DBConfigAccountStatusField = getStringConfig("database.zerotechhDB.UserMetaDataCollection.AccountStatusField", dbConfigManager)
)
