package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/smtp"

	"github.com/dipress/evostok/internal/sender"
	"github.com/pkg/errors"
)

func main() {

	var (
		host     = flag.String("host", "host", "enter smtp host")
		port     = flag.String("port", "port", "enter smpt port")
		password = flag.String("password", "password", "enter password for smtp server")
		login    = flag.String("login", "login", "enter login for e-vostok")
		secret   = flag.String("secret", "secret", "enter secret for e-vostok")
		id       = flag.String("id", "id", "service id for e-vostok")
	)

	flag.Parse()

	url := fmt.Sprintf(
		"https://clients.e-vostok.ru/checkbalance.php?login=%s&password=%s&service=%s",
		*login,
		*secret,
		*id,
	)
	mail := sender.Mail{}

	mail.SendlerID = "example@example.com"
	mail.ToIDs = []string{"example@example"}
	mail.Subject = "My subject"

	server := sender.SMTPServer{
		Host:     *host,
		Port:     *port,
		Password: *password,
	}

	client, err := conn(server.Host, server.Port, server.Password)

	if err != nil {
		log.Fatalf("open smtp client: %s", err)
	}

	defer client.Quit()

	send := sender.New(client, &server, url)

	if err := send.Send(&mail); err != nil {
		log.Fatal(err)
	}
}

func conn(host, port, password string) (*smtp.Client, error) {

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host + ":" + port,
	}

	conn, err := tls.Dial("tcp", host+":"+port, tlsconfig)

	if err != nil {
		return nil, errors.Wrapf(err, "Dial failed")
	}

	client, err := smtp.NewClient(conn, host+":"+port)

	if err != nil {
		return nil, errors.Wrapf(err, "Client not created")
	}

	return client, nil
}
