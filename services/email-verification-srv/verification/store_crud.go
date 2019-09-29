package verification

import (
	"time"

	"github.com/pkg/errors"

	"github.com/SonicRoshan/Velocity/global/config"
)

func (codeStore CodeStore) doesCodeWithIDExist(ID string) bool {
	var code config.VerificationCode
	filter := config.VerificationCode{
		ID: ID,
	}

	codeStore.mainCollection.FindOne(codeStore.ctx, filter).Decode(&code)

	return code != config.VerificationCode{}
}

//NewCode is used to generate new and store a verification code
func (codeStore CodeStore) NewCode(ID string) (string, error) {

	exists := codeStore.doesCodeWithIDExist(ID)
	if exists {
		return "", errors.New("Code With ID Already Exists")
	}

	codeStr := GenerateCode()

	code := config.VerificationCode{
		ID:          ID,
		CreationUTC: time.Now().Unix(),
		Code:        codeStr,
	}

	codeStore.mainCollection.InsertOne(codeStore.ctx, code)
	return code.Code, nil
}

//VerifyCode is used to verify a code and get the id tied to it
func (codeStore CodeStore) VerifyCode(codeStr string) string {
	var code config.VerificationCode
	filter := config.VerificationCode{
		Code: codeStr,
	}

	codeStore.mainCollection.FindOne(codeStore.ctx, filter).Decode(&code)

	return code.ID
}

//CleanUp is used to remove expired verification codes
func (codeStore CodeStore) CleanUp() {
	cursor, _ := codeStore.mainCollection.Find(codeStore.ctx, config.VerificationCode{})

	var code config.VerificationCode
	for cursor.Next(codeStore.ctx) {
		cursor.Decode(&code)
		expirationTime := time.Unix(code.CreationUTC, 0)
		if time.Now().After(expirationTime) {
			codeStore.mainCollection.DeleteOne(codeStore.ctx, code)
		}
	}
}
