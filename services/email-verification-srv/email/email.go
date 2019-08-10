package email

import (
	"crypto/tls"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

//SendSimpleEmail Is Used To Send A Basic Text Main
func SendSimpleEmail(text, toEmail string) error {
	dialer := gomail.NewDialer(
		config.SMTPAddress,
		config.SMTPPort,
		config.SMTPEmail,
		config.SMTPPassword,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: "smtp.gmail.com"}
	dialer.SSL = true
	message := gomail.NewMessage()

	message.SetBody("text", text)
	message.SetHeader("From", config.SMTPEmail)
	message.SetHeader("To", toEmail)

	err := dialer.DialAndSend(message)
	if err != nil && !config.DebugMode {
		return errors.Wrap(err, "Error while dialing and sending simple text mail")
	}
	return nil
}
