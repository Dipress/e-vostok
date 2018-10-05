package main

import (
	"crypto/tls"
	"net/smtp"

	"github.com/pkg/errors"
)

//AddTSL added TSL/SSL to server host
func AddTSL(host string) *tls.Config {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	return tlsconfig
}

//DialToServer func tls.dial to server
func DialToServer(serverName string, tslconfig *tls.Config) (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", serverName, tslconfig)

	if err != nil {
		return nil, errors.Wrapf(err, "Dial failed")
	}

	return conn, err
}

// CreateClient returns smtp client
func CreateClient(conn *tls.Conn, host string) (*smtp.Client, error) {
	client, err := smtp.NewClient(conn, host)

	if err != nil {
		return nil, errors.Wrapf(err, "Client not created")
	}

	return client, err
}
