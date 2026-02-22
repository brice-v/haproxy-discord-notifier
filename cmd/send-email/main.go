package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	from := "brice@example.com"
	to := []string{"other.brice@example.com"}
	addr := "127.0.0.1:9123"
	msg := []byte("From: brice@example.com\r\n" +
		"To: other-brice@example.com\r\n" +
		"Subject: Test mail\r\n\r\n" +
		"Email body\r\n")

	err := smtp.SendMail(addr, nil, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}
