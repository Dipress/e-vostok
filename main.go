package main

import (
	"log"
	"net/smtp"
)

func main() {
	balance := getBalance("https://clients.e-vostok.ru/checkbalance.php?login=krimea&password=Eefv8Dg4&service=21824")

	mail := Mail{}
	mail.SendlerID = "dmitry@crimeainfo.com"
	mail.ToIDs = []string{"tech@crimeainfo.com", "boleg@crimeainfo.com"}
	mail.Subject = "E-Vostok Balance"
	mail.Body = "The Balance is " + balance

	messageBody := mail.BuildMessage()

	smtpServer := SMTPServer{Host: "smtp.yandex.ru", Port: "465"}

	auth := smtp.PlainAuth("", mail.SendlerID, "dmitry2017", smtpServer.Host)

	tlsconfig := AddTSL(smtpServer.ServerName())

	conn, err := DialToServer(smtpServer.ServerName(), tlsconfig)

	if err != nil {
		log.Fatal(err)
	}
	client, err := CreateClient(conn, smtpServer.Host)

	if err != nil {
		log.Fatal(err)
	}

	// User Auth
	if err := client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// Add all FROM and TO
	if err := client.Mail(mail.SendlerID); err != nil {
		log.Panic(err)
	}

	for _, k := range mail.ToIDs {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// DATA
	w, err := client.Data()

	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}
