package config

import "github.com/dgrijalva/jwt-go"

//UserType is used to store user data
type UserType struct {
	UserID             string `bson:"id,omitempty,-"`
	Username           string `bson:"Username,omitempty,-"`
	Email              string `bson:"Email,omitempty,-"`
	Password           string `bson:"Password,omitempty,-"`
	FirstName          string `bson:"FirstName,omitempty,-"`
	LastName           string `bson:"LastName,omitempty,-"`
	Gender             string `bson:"Gender,omitempty,-"`
	BirthdayUTC        int64  `bson:"BirthdayUTC,omitempty,-"`
	AccountStatus      string `bson:"AccountStatus,omitempty,-"`
	AccountCreationUTC int64  `bson:"AccountCreationUTC,omitempty,-"`
}

//JWTClaims is used to store jwt claims
type JWTClaims struct {
	UserIdentity  string   `json:"UserIdentity"`
	IsFresh       bool     `json:"IsFresh"`
	IsRefresh     bool     `json:"IsRefresh"`
	Scopes        []string `json:"Scopes"`
	CreationUTC   int64    `json:"CreationUTC"`
	ExpirationUTC int64    `json:"ExpirationUTC"`
	jwt.StandardClaims
}
