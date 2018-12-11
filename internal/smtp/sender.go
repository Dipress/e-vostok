package smtp

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/pkg/errors"
)

const subject = "My subject"

type Sender struct {
	client *smtp.Client
	from   string
	to     []string
}

// New is factory function for basic struct
func New(client *smtp.Client, from string) *Sender {
	s := Sender{
		client: client,
		from:   from,
	}

	return &s
}

func (s *Sender) Send(body string, to []string) error {
	if err := s.client.Mail(s.from); err != nil {
		return errors.Wrap(err, "mail transaction initiate failed")
	}

	for _, k := range to {
		if err := s.client.Rcpt(k); err != nil {
			return errors.Wrap(err, "rcpt call failed")
		}
	}

	w, err := s.client.Data()

	if err != nil {
		return errors.Wrap(err, "data command failed")
	}

	mail := buildMessage(body, s.from, subject, to)

	_, err = w.Write(mail)
	if err != nil {
		return errors.Wrap(err, "write command failed")
	}

	if err = w.Close(); err != nil {
		return errors.Wrap(err, "writer close failed")
	}

	return nil
}

func buildMessage(body, subject, from string, to []string) []byte {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "From: %s\r\n", from)

	if len(to) > 0 {
		fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(to, ";"))
	}

	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprint(&buf, "\r\n"+body)

	return buf.Bytes()
}
