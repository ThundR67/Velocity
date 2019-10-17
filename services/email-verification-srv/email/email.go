package email

import (
	"net/smtp"
	"strconv"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/pkg/errors"
)

//SendSimpleEmail Is Used To Send A Basic Text Main
func SendSimpleEmail(text, toEmail string) error {

	auth := smtp.PlainAuth(
		"",
		config.SMTPEmail,
		config.SMTPPassword,
		config.SMTPAddress,
	)

	err := smtp.SendMail(
		config.SMTPAddress+":"+strconv.Itoa(config.SMTPPort),
		auth,
		config.SMTPEmail,
		[]string{toEmail},
		[]byte(text),
	)

	if err != nil && (!config.DebugMode || config.VerificationSendEmail) {
		return errors.Wrap(err, "Error while sending email")
	}

	return nil
}
