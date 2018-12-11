package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/smtp"

	"github.com/dipress/evostok/internal/http"
	"github.com/dipress/evostok/internal/send"
	smtpCli "github.com/dipress/evostok/internal/smtp"
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

	sendlerID := "example@example.com"
	to := []string{"example@example.com"}

	client, err := conn(*host, *port, *password)
	if err != nil {
		log.Fatalf("open smtp client: %s", err)
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", sendlerID, *password, *host+":"+*port)

	if err := client.Auth(auth); err != nil {
		log.Fatalf("smtp authenticate failed: %v", err)
	}

	g := http.NewBalance()
	s := smtpCli.New(client, sendlerID)

	srv := send.NewService(g, s)
	if err := srv.Deliver(url, to); err != nil {
		log.Fatalf("failed to send: %v", err)
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
