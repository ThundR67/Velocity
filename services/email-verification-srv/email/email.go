package email

import (
	"crypto/tls"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

//SendSimpleEmail Is Used To Send A Basic Text Main
func SendSimpleEmail(text string) error {
	dialer := gomail.NewDialer(
		config.SMTPAddress,
		config.SMTPPort,
		config.SMTPEmail,
		config.SMTPPassword,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	message := gomail.NewMessage()
	message.SetBody("text", text)
	err := dialer.DialAndSend(message)
	if err != nil {
		return errors.Wrap(err, "Error while dialing and sending simple text mail")
	}
	return nil
}
