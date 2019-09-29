package verification

import (
	"github.com/google/uuid"
)

//GenerateCode is used to generate a code for email verification
func GenerateCode() string {
	code, _ := uuid.NewRandom()
	return code.String()
}
