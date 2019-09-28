package verification

import (
	"time"

	
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/global/config"
)

func (codeStore CodeStore) doesCodeWithIDExist(ID string) (bool, error) {
	var code config.VerificationCode
	filter := config.VerificationCode{
		ID: ID,
	}

	err := codeStore.mainCollection.FindOne(codeStore.ctx, filter).Decode(&code)
	if err != nil {
		return false, nil
	}

	return code != config.VerificationCode{}, nil
}

//NewCode is used to generate new and store a verification code
func (codeStore CodeStore) NewCode(ID string) (string, error) {
	exists, err := codeStore.doesCodeWithIDExist(ID)
	if err != nil {
		return "", errors.Wrap(err, "Error while checking if code with id exists")
	} else if exists {
		return "", errors.New("Code With ID Already Exists")
	}

	codeStr, err := GenerateCode()
	if err != nil {
		return "", errors.Wrap(err, "Error while genrating code")
	}

	code := config.VerificationCode{
		ID:          ID,
		CreationUTC: time.Now().Unix(),
		Code:        codeStr,
	}

	_, err = codeStore.mainCollection.InsertOne(codeStore.ctx, code)
	if err != nil {
		return "", errors.Wrap(err, "Error while inserting code into DB")
	}
	return code.Code, nil
}

//VerifyCode is used to verify a code and get the id tied to it
func (codeStore CodeStore) VerifyCode(codeStr string) (string, error) {
	var code config.VerificationCode
	filter := config.VerificationCode{
		Code: codeStr,
	}

	err := codeStore.mainCollection.FindOne(codeStore.ctx, filter).Decode(&code)
	if err != nil {
		err = errors.Wrap(err, "Error while finding code with code string")
		return "", err
	}

	return code.ID, nil
}

//CleanUp is used to remove expired verification codes
func (codeStore CodeStore) CleanUp() {
	cursor, err := codeStore.mainCollection.Find(codeStore.ctx, config.VerificationCode{})
	if err != nil {
		log.Error("Error while clearing up expired verification code", zap.Error(err))
	}

	var code config.VerificationCode
	for cursor.Next(codeStore.ctx) {
		cursor.Decode(&code)
		expirationTime := time.Unix(code.CreationUTC, 0)
		if time.Now().After(expirationTime) {
			codeStore.mainCollection.DeleteOne(codeStore.ctx, code)
		}
	}
}
