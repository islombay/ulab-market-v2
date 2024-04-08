package smtp

import (
	"app/config"
	"app/pkg/helper"
	"app/pkg/logs"
	"fmt"
	"net/smtp"
	"os"
)

type emailSMTPService struct {
	cfg *config.SMTPEmailConfig
	log logs.LoggerInterface
}

func NewEmailSMTPService(cfg *config.SMTPEmailConfig, log logs.LoggerInterface) *emailSMTPService {
	return &emailSMTPService{
		cfg: cfg,
		log: log,
	}
}

func (se *emailSMTPService) SendVerificationCode(to, code string) error {
	if !helper.IsValidEmail(to) {
		return helper.ErrInvalidEmail
	}
	msg := fmt.Sprintf("Assalamu Alaykum\nTasdiqlash kodingiz: %s", code)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", se.cfg.Host, se.cfg.Port),
		smtp.PlainAuth("", se.cfg.SenderEmail, os.Getenv("SMTP_EMAIL_PWD"), se.cfg.Host),
		se.cfg.SenderEmail,
		[]string{to},
		[]byte(msg),
	)
	if err != nil {
		se.log.Error("could not send email", logs.Error(err), logs.String("email", to))
		return err
	}
	return nil
}
