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
	userID, msg, err := client.Add(mockUserData, mockUserExtraData)
	assert.NoError(err, "Adding User Returned Error")
	assert.Equal("", msg, "Message By Add Should Be Blank")

	//Testing Get
	data, msg, err := client.Get(userID)
	assert.NoError(err, "Getting User Returned Error")
	assert.Equal("", msg, "Message By Get Should Be Blank")
	assert.Equal(mockUserData.Username, data.Username, "Username Should Match")

	//Testing Update
	updatedEmail := "sonicroshan122@gmail.com"
	update := config.UserMain{
		Email: updatedEmail,
	}
	err = client.Update(userID, update)
	assert.NoError(err, "UpdateUser Returned Error")
	data, _, _ = client.Get(userID)
	assert.Equal(updatedEmail, data.Email, "Updated Email Should Match Email In DB")

	//Testing Auth
	id, msg, err := client.Auth(mockUserData.Username, mockUserData.Password)
	assert.NoError(err, "Auth User Returned Error")
	assert.Equal("", msg, "Message By Auth Should Be Blank")
	assert.Equal(userID, id, "Ids should match")

	msg, err = client.Activate(updatedEmail)
	assert.NoError(err)
	assert.Zero(msg)

}
