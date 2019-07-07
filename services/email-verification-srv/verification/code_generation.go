package verification

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

//GenerateCode is used to generate a code for email verification
func GenerateCode() (string, error) {
	code, err := uuid.NewRandom()
	if err != nil {
		err = errors.Wrap(err, "Error While Generating UUID")
	}
	return code.String(), err
}
