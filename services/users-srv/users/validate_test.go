package users

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/utils"
	"github.com/stretchr/testify/assert"
)

func validateMain(key, value string) bool {
	toValidateStr := fmt.Sprintf(`{"%s":"%s"}`, key, value)
	toValidate, _ := utils.GetMockUserData()
	err := json.Unmarshal([]byte(toValidateStr), &toValidate)
	if err != nil {
		panic(fmt.Sprintf("Key %s Val %s Err %s", key, value, err.Error()))
	}
	return validateUserMainData(toValidate)
}

func validateExtra(key, value string) bool {
	toValidateStr := fmt.Sprintf(`{"%s":"%s"}`, key, value)
	_, toValidate := utils.GetMockUserData()
	err := json.Unmarshal([]byte(toValidateStr), &toValidate)
	if err != nil {
		panic(fmt.Sprintf("Key %s Val %s Err %s", key, value, err.Error()))
	}
	return validateUserExtraData(toValidate)
}

func TestGenderValidator(t *testing.T) {
	assert := assert.New(t)
	for _, gender := range config.UserDataConfigGenderTypes {
		assert.True(isGender(gender))
	}
	assert.False(isGender("invalidgender"))
}

func TestValidators(t *testing.T) {
	assert := assert.New(t)

	mainData := map[string]map[string][]string{
		"Username": map[string][]string{
			"valid":   []string{"testing", "test123"},
			"invalid": []string{"ta", "", "Atest", " test"},
		},
		"Email": map[string][]string{
			"valid":   []string{"sonicroshan122@gmail.com", "vishvam.shashtri@gmail.com"},
			"invalid": []string{"test", "", "fakeEmail", " sonicroshan122@gmail.com "},
		},
		"Password": map[string][]string{
			"valid":   []string{"testing", "test123"},
			"invalid": []string{"test", "", "123", " 1"},
		},
	}

	extraData := map[string]map[string][]string{
		"FirstName": map[string][]string{
			"valid":   []string{"Roshan"},
			"invalid": []string{" ", "A", " Roshan"},
		},
		"LastName": map[string][]string{
			"valid":   []string{"Roshan"},
			"invalid": []string{" ", "A", " Roshan"},
		},
		"Gender": map[string][]string{
			"valid":   []string{"male"},
			"invalid": []string{"invalidgender"},
		},
	}

	for key, testCase := range mainData {
		for _, valid := range testCase["valid"] {
			assert.True(validateMain(key, valid))
		}

		for _, invalid := range testCase["invalid"] {
			assert.False(validateMain(key, invalid))
		}
	}

	for key, testCase := range extraData {
		for _, valid := range testCase["valid"] {
			assert.True(validateExtra(key, valid))
		}

		for _, invalid := range testCase["invalid"] {
			assert.False(validateExtra(key, invalid))
		}
	}

	//Testing validation for birthday utc
	_, userExtra := utils.GetMockUserData()
	assert.True(validateUserExtraData(userExtra))
	userExtra.BirthdayUTC = time.Now().Unix()
	assert.False(validateUserExtraData(userExtra))

}
