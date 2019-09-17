package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
)

func generateRandomString(length int) string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = config.UserDataConfigUserIDCharset[seededRand.Intn(len(config.UserDataConfigUserIDCharset))]
	}
	return string(b)
}

//GetMockUserData returns mock user data for testing
func GetMockUserData() (config.UserMain, config.UserExtra) {
	mockUsername := strings.ToLower(generateRandomString(10))
	mockPassword := generateRandomString(10)
	mockUserData := config.UserMain{
		Username: mockUsername,
		Password: mockPassword,
		Email:    mockUsername + "@gmail.com",
	}

	mockUserExtraData := config.UserExtra{
		FirstName:   mockUsername,
		LastName:    mockUsername,
		Gender:      config.UserDataConfigGenderTypeMale,
		BirthdayUTC: int64(864466669),
	}

	return mockUserData, mockUserExtraData
}
