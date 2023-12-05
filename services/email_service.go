package services

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailService(config *cfg.Config) EmailService {
	user := config.SMTPUser
	password := config.SMTPPassword
	host := config.SMTPHost
	port := config.SMTPPort

	d := gomail.NewDialer(host, port, user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return EmailService{
		dialer: d,
		from:   config.EmailFrom,
	}
}

func (service *EmailService) Send(user *models.User, data *models.Email, temp *template.Template, templateName string) error {
	to := user.Email

	var body bytes.Buffer

	if err := temp.ExecuteTemplate(&body, templateName, &data); err != nil {
		log.Fatal("Could not execute template", err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", service.from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", data.Subject)
	message.SetBody("text/html", body.String())
	message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	if err := service.dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
