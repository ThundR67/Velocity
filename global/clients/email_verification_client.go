package clients

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	proto "github.com/SonicRoshan/Velocity/services/email-verification-srv/proto"
	micro "github.com/micro/go-micro"
	"github.com/pkg/errors"
)

//NewEmailVerificationClient is used to make a email verification Client
func NewEmailVerificationClient(service micro.Service) EmailVerificationClient {
	client := EmailVerificationClient{}
	client.Init(service)
	return client
}

//EmailVerificationClient is jwt service client
type EmailVerificationClient struct {
	client proto.EmailVerificationService
}

//Init initializes cliet
func (emailVerificationClient *EmailVerificationClient) Init(service micro.Service) {
	emailVerificationClient.client = proto.NewEmailVerificationService(config.EmailVerificationSrv, service.Client())
}

//SendVerification is used to send verification email to a user
func (emailVerificationClient EmailVerificationClient) SendVerification(email string) error {
	request := proto.SendVerificationRequest{
		Email: email,
	}

	_, err := emailVerificationClient.client.SendVerification(context.TODO(), &request)
	if err != nil {
		return errors.Wrap(err, "Error While Sending Verification Through Client")
	}

	return nil
}

//Verify is used to verify verification code
func (emailVerificationClient EmailVerificationClient) Verify(code string) (string, error) {
	request := proto.VerifyRequest{
		VerificationCode: code,
	}

	response, err := emailVerificationClient.client.Verify(context.TODO(), &request)
	if err != nil {
		return "", errors.Wrap(err, "Error While Verifying Code Through Client")
	}

	return response.Email, nil
}
