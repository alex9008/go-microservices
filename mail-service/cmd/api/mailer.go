package main

import (
	"bytes"
	"html/template"
	"log"
	"strconv"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        string
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {

	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formatterMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTestMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port, _ = strconv.Atoi(m.Port)
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formatterMessage)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (m *Mail) buildPlainTestMessage(msg Message) (string, error) {

	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	if err = t.ExecuteTemplate(&buffer, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := buffer.String()

	return plainMessage, nil

}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	if err = t.ExecuteTemplate(&buffer, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := buffer.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil

}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(encryption string) mail.Encryption {
	switch encryption {
	case "ssl":
		return mail.EncryptionSSL
	case "tls":
		return mail.EncryptionTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
