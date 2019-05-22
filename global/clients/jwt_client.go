package clients

import (
	"context"

	proto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	micro "github.com/micro/go-micro"
	"github.com/pkg/errors"
)

//JWTClient is jwt service client
type JWTClient struct {
	client proto.JWTManagerService
}

//Init initalizes cliet
func (jwtClient *JWTClient) Init(service micro.Service) {
	jwtClient.client = proto.NewJWTManagerService("jwt-srv", service.Client())
}

//GenerateFreshAccessToken generates fresh access token
func (jwtClient JWTClient) GenerateFreshAccessToken(userIdentity string) (string, error) {
	request := proto.JWTData{
		UserIdentity: userIdentity,
	}
	response, err := jwtClient.client.GenerateFreshAccessToken(context.TODO(), &request)
	if err != nil {
		return "", errors.Wrap(err, "Error While Generating Fresh Access Token Through Client Through Client")
	}
	return response.Token, err
}

//GenerateAccessAndRefreshToken generates access and refresh token based on userIdentity
func (jwtClient JWTClient) GenerateAccessAndRefreshToken(userIdentity string, scopes []string) (string, string, string, error) {
	request := proto.JWTData{
		UserIdentity: userIdentity,
		Scopes:       scopes,
	}
	response, err := jwtClient.client.GenerateAccessAndRefreshToken(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error While Generating Access And Refresh Token Through Client")
		return "", "", "", err
	}
	return response.AcccessToken, response.RefreshToken, response.Message, err
}

//GenerateAccessAndRefreshTokenBasedOnRefreshToken generates access and refresh token based on previous refresh token
func (jwtClient JWTClient) GenerateAccessAndRefreshTokenBasedOnRefreshToken(refreshToken string) (string, string, string, error) {
	request := proto.Token{
		Token: refreshToken,
	}
	response, err := jwtClient.client.GenerateAccessAndRefreshTokenBasedOnRefreshToken(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error While Generating Access And Refresh Token Bases On Refresh Token Through Client")
		return "", "", "", err
	}
	return response.AcccessToken, response.RefreshToken, response.Message, err
}

//ValidateFreshAccessToken validates fresh access token
func (jwtClient JWTClient) ValidateFreshAccessToken(freshAccessToken string) (bool, string, error) {
	request := proto.Token{
		Token: freshAccessToken,
	}
	response, err := jwtClient.client.ValidateFreshAccessToken(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error While Validating Fresh Access Token Through Client")
		return false, "", err
	}
	return response.Valid, response.Message, err
}

//ValidateToken validates access token and returnes its claims
func (jwtClient JWTClient) ValidateToken(token string) (bool, string, []string, error) {
	request := proto.Token{
		Token: token,
	}
	response, err := jwtClient.client.ValidateToken(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error While Validating Token Through Client")
		return false, "", nil, err
	}
	userIdentity := response.JwtData.UserIdentity
	scopes := response.JwtData.Scopes
	return response.Valid, userIdentity, scopes, err
}
