package email

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"gopkg.in/gomail.v2"
)

type IEmailService interface {
	Send(to string, subject string, message string) error
}

type EmailService struct {
	cfg *config.SMTPConfig
}

func NewEmailService(cfg *config.SMTPConfig) IEmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) Send(to string, subject string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.ADDRESS)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := gomail.NewDialer(s.cfg.HOST, s.cfg.PORT, s.cfg.ADDRESS, s.cfg.PASSWORD)

	return d.DialAndSend(m)
}
