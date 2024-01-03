package email

import (
	"auth/internal/config"
	"fmt"
	"net/smtp"
)

func SendVerificationCodeEmail(email, verificationCode string, config config.EmailConfig) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Verification Code\n\nYour verification code is: %s", config.From, email, verificationCode)

	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)

	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.From, []string{email}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
