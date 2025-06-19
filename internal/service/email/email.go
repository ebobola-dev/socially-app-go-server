package email

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"gopkg.in/gomail.v2"
)

type IEmailService interface {
	Send(to string, subject string, message string) error
}

type emailService struct {
	cfg *config.SMTPConfig
}

func NewEmailService(cfg *config.SMTPConfig) IEmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) Send(to string, subject string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.ADDRESS)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := gomail.NewDialer(s.cfg.HOST, s.cfg.PORT, s.cfg.ADDRESS, s.cfg.PASSWORD)

	return d.DialAndSend(m)
}
