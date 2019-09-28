package handler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
	"github.com/SonicRoshan/Velocity/services/email-verification-srv/email"
	proto "github.com/SonicRoshan/Velocity/services/email-verification-srv/proto"
	"github.com/SonicRoshan/Velocity/services/email-verification-srv/verification"
	
)

var log = logger.GetLogger("email_verification_service.log")

//EmailVerification is used to handle for email verification service
type EmailVerification struct {
	codeStore verification.CodeStore
}

//Init is used to initialize the handler
func (emailVerification *EmailVerification) Init() error {
	emailVerification.codeStore = verification.CodeStore{}
	err := emailVerification.codeStore.Init()
	if err != nil {
		log.Error(
			"Error While Initializing Code Store",
			zap.Error(err),
		)
		return errors.Wrap(err, "Error While Initializing Code Store")
	}

	//Setting up cleaner
	go func() {
		for {
			time.Sleep(config.VerificationCleanerSleepTime)
			emailVerification.codeStore.CleanUp()
		}
	}()

	return nil
}

//SendVerification is used to handler SendVerification function
func (emailVerification EmailVerification) SendVerification(
	ctx context.Context,
	request *proto.SendVerificationRequest,
	response *proto.SendVerificationResponse) error {

	log.Debug(
		"Sending A Verification",
		zap.String("Email", request.Email),
	)

	code, err := emailVerification.codeStore.NewCode(request.Email)
	if err != nil {
		msg := "Error While Generating Verification Code"
		log.Error(msg, zap.Error(err))
		return errors.Wrap(err, msg)
	}

	log.Info(
		"Created A Code",
		zap.String("Code", code),
	)

	err = email.SendSimpleEmail(code, request.Email)
	if err != nil {
		msg := "Error While Emailing Verification Code"
		log.Error(msg, zap.Error(err))
		return errors.Wrap(err, msg)
	}

	log.Info(
		"Sent A Code",
		zap.String("Code", code),
		zap.String("Email", request.Email),
	)
	return nil
}

//Verify is used to handle Verify Function
func (emailVerification EmailVerification) Verify(
	ctx context.Context,
	request *proto.VerifyRequest,
	response *proto.VerifyResponse) error {

	log.Debug("Verifying Code", zap.String("Code", request.VerificationCode))

	email, err := emailVerification.codeStore.VerifyCode(request.VerificationCode)
	if err != nil {
		msg := "Error While Verifying Code"
		log.Error(msg, zap.Error(err))
		return errors.Wrap(err, msg)
	}

	response.Email = email
	log.Info("Verified Email")
	return nil
}
