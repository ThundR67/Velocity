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
	//DBConfigMainDB Is The ZeroTechhDB Name
	DBConfigMainDB = getStringConfig("database.zerotechhDB", dbConfigManager)

	//DBConfigUserMainDataCollection Is UserData Collection
	DBConfigUserMainDataCollection = getStringConfig("database.zerotechhDB.UserDataCollection", dbConfigManager)
	//DBConfigUserExtraDataCollection Is Where Extra Data Is Stored
	DBConfigUserExtraDataCollection = getStringConfig("database.zerotechhDB.UserExtraDataCollection", dbConfigManager)
	//DBConfigUserMetaDataCollection Is Where User's meta Data Is Stored
	DBConfigUserMetaDataCollection = getStringConfig("database.zerotechhDB.UserMetaDataCollection", dbConfigManager)
)
