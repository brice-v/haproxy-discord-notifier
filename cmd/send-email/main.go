package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {

	from := "john.doe@example.com"

	user := "9c1d45eaf7af5b"
	password := "ad62926fa75d0f"

	to := []string{
		"roger.roe@example.com",
	}

	addr := "127.0.0.1:9123"
	host := "local-dev"

	msg := []byte("From: john.doe@example.com\r\n" +
		"To: roger.roe@example.com\r\n" +
		"Subject: Test mail\r\n\r\n" +
		"Email body\r\n")

	auth := smtp.PlainAuth("", user, password, host)

	err := smtp.SendMail(addr, auth, from, to, msg)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}
