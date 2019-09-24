package config

import "github.com/dgrijalva/jwt-go"

//UserMain is used to store user main data
type UserMain struct {
	UserID   string `bson:"_id,omitempty,-" isArg:"true" description:"Unique User Identity" readScope:"user:read:user_id"`
	Username string `bson:"Username,omitempty,-" isArg:"true" description:"Username of the user" readScope:"user:read:username"`
	Email    string `bson:"Email,omitempty,-" isArg:"true" description:"Email of the user" readScope:"user:read:email"`
	Password string `bson:"Password,omitempty,-" exclude:"true"`
}

//UserExtra is used to store user extra data
type UserExtra struct {
	UserID      string `bson:"_id,omitempty,-" readScope:"user:read:user_id"`
	FirstName   string `bson:"FirstName,omitempty,-" readScope:"user:read:first_name"`
	LastName    string `bson:"LastName,omitempty,-" readScope:"user:read:last_name"`
	Gender      string `bson:"Gender,omitempty,-" readScope:"user:read:gender"`
	BirthdayUTC int64  `bson:"BirthdayUTC,omitempty,-" readScope:"user:read:birthday_utc"`
}

//UserMeta is used to store user meta data
type UserMeta struct {
	UserID             string `bson:"_id,omitempty,-"`
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
