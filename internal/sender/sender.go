package sender

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/pkg/errors"
)

// Mail struct
type Mail struct {
	SendlerID string
	ToIDs     []string
	Subject   string
	Body      string
}

// SMTPServer struct for server settings
type SMTPServer struct {
	Host     string
	Port     string
	Password string
}

type basic struct {
	client *smtp.Client
	server *SMTPServer
	url    string
}

//Sender Interface for send email messages
type Sender interface {
	Send(mail *Mail) error
}

//New is factory function for basic struct
func New(client *smtp.Client, server *SMTPServer, url string) Sender {
	b := basic{
		client: client,
		server: server,
		url:    url,
	}

	return &b
}

func (b *basic) Send(mail *Mail) error {
	res, err := http.Get(b.url)

	if err != nil {
		errors.Wrap(err, "get request failed")
	}

	result, err := ioutil.ReadAll(res.Body)

	if err != nil {
		errors.Wrap(err, "collect data failed")
	}

	defer res.Body.Close()

	mail.Body = "The Balance is " + string(result)

	auth := smtp.PlainAuth("", mail.SendlerID, b.server.Password, b.server.Host+":"+b.server.Port)

	if err := b.client.Auth(auth); err != nil {
		return errors.Wrap(err, "smtp authenticate failed")
	}

	if err := b.client.Mail(mail.SendlerID); err != nil {
		return errors.Wrap(err, "mail transaction initiate failed")
	}

	for _, k := range mail.ToIDs {
		if err := b.client.Rcpt(k); err != nil {
			return errors.Wrap(err, "rcpt call failed")
		}
	}

	w, err := b.client.Data()

	if err != nil {
		return errors.Wrap(err, "data command failed")
	}

	_, err = w.Write([]byte(mail.buildMessage()))

	if err != nil {
		return errors.Wrap(err, "write command failed")
	}

	err = w.Close()

	if err != nil {
		return errors.Wrap(err, "writer close failed")
	}

	return nil
}

func (mail *Mail) buildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.SendlerID)

	if len(mail.ToIDs) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.ToIDs, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	message += "\r\n" + mail.Body

	return message
}
