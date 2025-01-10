package email

import (
	"bytes"
	"fmt"
	"net/smtp"

	"github.com/Vaibhavsahu2810/redis-go/internal/config"
	"github.com/Vaibhavsahu2810/redis-go/internal/templates"
)

type Sender struct {
	config    *config.Config
	templates *templates.Manager
}

func NewSender(cfg *config.Config, tmpl *templates.Manager) *Sender {
	return &Sender{
		config:    cfg,
		templates: tmpl,
	}
}

func (s *Sender) SendEmail(to, subject, templateName string, data map[string]interface{}) error {
	body, err := s.templates.RenderWithSafeURLs(templateName, data)
	if err != nil {
		return fmt.Errorf("error rendering template: %w", err)
	}

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	var message bytes.Buffer
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.EmailFromName, s.config.EmailFrom))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	message.WriteString(body)

	// Send the email
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	if err := smtp.SendMail(addr, auth, s.config.EmailFrom, []string{to}, message.Bytes()); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
