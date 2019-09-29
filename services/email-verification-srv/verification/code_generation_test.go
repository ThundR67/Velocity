package verification

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestCodeGeneration(t *testing.T) {
	assert := assert.New(t)

	code := GenerateCode()
	assert.NotZero(code)
	assert.True(govalidator.IsUUIDv4(code))

	codeStore := CodeStore{}
	codeStore.Init()
	codeStore.CleanUp()

	code, err := codeStore.NewCode("123")
	assert.NoError(err)
	assert.NotZero(code)
	assert.True(govalidator.IsUUIDv4(code))

	id := codeStore.VerifyCode(code)
	assert.Equal("123", id)

	id = codeStore.VerifyCode("1")
	assert.Zero(id)

	code, err = codeStore.NewCode("123")
	assert.Error(err)
	assert.Zero(code)

	exists := codeStore.doesCodeWithIDExist("123")
	assert.True(exists)
	exists = codeStore.doesCodeWithIDExist("1")
	assert.False(exists)

}
