package clients

import (
	"testing"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/utils"
	micro "github.com/micro/go-micro"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	assert := assert.New(t)
	assert.Error(copy(32, 45))
}

func TestUsersClient(t *testing.T) {
	assert := assert.New(t)
	service := micro.NewService(micro.Name("test"))
	client := NewUsersClient(service)

	mockUserData, mockUserExtraData := utils.GetMockUserData()

	//Testing Add
	userID, msg := client.Add(mockUserData, mockUserExtraData)
	assert.Zero(msg)

	//Testing Get
	data, err := client.Get("InvalidUserID")
	assert.Zero(data)
	assert.Error(err)

	data, err = client.Get(userID)
	assert.NoError(err)
	assert.Equal(mockUserData.Username, data.Username)
	assert.NotZero(data)

	data, msg = client.GetByUsernameOrEmail(mockUserData.Username, "")
	assert.Zero(msg)
	assert.NotZero(data)

	data, msg = client.GetByUsernameOrEmail("", mockUserData.Email)
	assert.Zero(msg)
	assert.NotZero(data)

	//Testing GetExtra
	extra, err := client.GetExtra("InvalidUserID")
	assert.Zero(extra)
	assert.Error(err)

	extra, err = client.GetExtra(userID)
	assert.NoError(err)
	assert.Equal(mockUserExtraData.LastName, extra.LastName)
	assert.NotZero(extra)

	//Testing Update
	updatedEmail := "sonicroshan122@gmail.com"
	update := config.UserMain{
		Email: updatedEmail,
	}
	client.Update(userID, update)
	data, _ = client.Get(userID)
	assert.Equal(updatedEmail, data.Email)

	updatedFirstName := "NewName"
	updateExtra := config.UserExtra{
		FirstName: updatedFirstName,
	}
	client.UpdateExtra(userID, updateExtra)
	extra, _ = client.GetExtra(userID)
	assert.Equal(updatedFirstName, extra.FirstName)

	//Testing Auth
	id, msg, err := client.Auth(mockUserData.Username, mockUserData.Password)
	assert.NoError(err)
	assert.Equal("", msg)
	assert.Equal(userID, id)

	msg = client.Activate(updatedEmail)
	assert.Zero(msg)

}
