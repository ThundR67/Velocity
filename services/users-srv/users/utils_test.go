package users

import (
	"testing"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert := assert.New(t)

	metaData := generateUserMetaData()
	assert.Equal(config.UserDataConfigAccountStatusUnactivated,
		metaData.AccountStatus,
		"GenerateMetaData Should Return Data With Account Status Unactivated")
	assert.WithinDuration(time.Unix(metaData.AccountCreationUTC, 0),
		time.Now(),
		time.Second*1,
		"Account Creation Time Is Incorrect",
	)
}
