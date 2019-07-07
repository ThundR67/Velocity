package verification

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestCodeGeneration(t *testing.T) {
	assert := assert.New(t)

	code, err := GenerateCode()
	assert.NoError(err, "GenerateCode Returned Error")
	assert.NotEqual("", code, "Generated Code Should Not Be Empty")
	assert.True(govalidator.IsUUIDv4(code), "Code Must Be Valid UUID v4")

}
