package main

import (
	"fmt"
	"strings"
)

// Mail struct
type Mail struct {
	SendlerID string
	ToIDs     []string
	Subject   string
	Body      string
}

// BuildMessage makes simple mail object
func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.SendlerID)
	if len(mail.ToIDs) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.ToIDs, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	message += "\r\n" + mail.Body

	return message
}
