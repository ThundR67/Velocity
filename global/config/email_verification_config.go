package config

import (
	"time"
)

var verificationConfig = getConfigManager("email_verification.config")

var (
	//VerificationMainCollection is the main collection where data will be stored
	VerificationMainCollection = getStringConfig("verificationCollection", verificationConfig)
	//VerificationExpirationTimeMinutes is time in which email verification code will expire
	VerificationExpirationTimeMinutes = time.Minute * time.Duration(getIntConfig("verificationExpirationTimeMinutes", verificationConfig))
	//VerificationCleanerSleepTime is time after which cleaner runs
	VerificationCleanerSleepTime = time.Hour * time.Duration(getIntConfig("verificationCleanerSleepHour", verificationConfig))

	//SMTPAddress is the smtp server address
	SMTPAddress = getStringConfig("smtp.Address", verificationConfig)
	//SMTPPort is the smtp server port
	SMTPPort = getIntConfig("smtp.Port", verificationConfig)
	//SMTPEmail is the smtp server's email
	SMTPEmail = getStringConfig("smtp.Email", verificationConfig)
	//SMTPPassword is the smtp server's password
	SMTPPassword = getStringConfig("smtp.Password", verificationConfig)
)

const (
	/*VerificationSendEmail determines if server will send emails or not.
	However, if debufMode is False,
	then emails will be sent no matter if this setting is false*/
	VerificationSendEmail = false
)

//VerificationCode is used to store verification data into db
type VerificationCode struct {
	ID          string `bson:"_id,omitempty,-"`
	Code        string `bson:"code,omitempty,-"`
	CreationUTC int64  `bson:"creationUTC,omitempty,-"`
}
