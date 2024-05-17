package smtp

import (
	"app/config"
	"app/pkg/logs"
)

type smtpService struct {
	log logs.LoggerInterface
	cfg *config.SMTPConfig

	email EmailSMTPInterface
	phone PhoneSMTPInterface
}

type EmailSMTPInterface interface {
	SendVerificationCode(to, code string) error
}

type PhoneSMTPInterface interface {
	SendVerificationCode(to, code string) error
}

type SMTPInterface interface {
	Email() EmailSMTPInterface
	Phone() PhoneSMTPInterface
}

func NewSMTPService(log logs.LoggerInterface, cfg *config.SMTPConfig) SMTPInterface {
	return &smtpService{
		log:   log,
		cfg:   cfg,
		email: NewEmailSMTPService(&cfg.Email, log),
		phone: NewPhoneSMTPService(&cfg.Email, log),
	}
}

func (s *smtpService) Email() EmailSMTPInterface {
	if s.email == nil {
		s.email = NewEmailSMTPService(&s.cfg.Email, s.log)
	}
	return s.email
}

func (s *smtpService) Phone() PhoneSMTPInterface {
	if s.phone == nil {
		s.phone = NewPhoneSMTPService(&s.cfg.Email, s.log)
	}
	return s.phone
}
