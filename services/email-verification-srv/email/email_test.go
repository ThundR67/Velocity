package email

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SonicRoshan/Velocity/global/config"
)

func TestEmail(t *testing.T) {
	assert := assert.New(t)
	assert.Contains(config.SMTPEmail, "@", "In Email Config, Valid SMTP Settings Have Not Been Added")
	err := SendSimpleEmail("TEST", config.SMTPEmail)
	assert.NoError(err)

	err = SendSimpleEmail("TEST", "Invalid")
	assert.Error(err)
}
