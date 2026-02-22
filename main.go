package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

const (
	webhookUrlEnvVariable = "WEBHOOK_URL"
	expectedUrlPrefix     = "https://discord.com/api/webhooks/"
)

var url = ""

func verifyUrlShape() {
	if url == "" {
		log.Fatalf("Provide `%s` through environment variables. Got empty string", webhookUrlEnvVariable)
	}
	if !strings.HasPrefix(url, expectedUrlPrefix) {
		log.Fatalf("URL `%s` does not look like a discord webhook url. should start with %s\n", url, expectedUrlPrefix)
	}
	urlParts := strings.Split(url, expectedUrlPrefix)
	if len(urlParts) != 2 {
		log.Fatalf("After splitting the url on `%s` 2 parts were not found. The url should have /{webhook.id}/{webhook.token} at the end\n", expectedUrlPrefix)
	}
	finalParts := strings.Split(urlParts[1], "/")
	if len(finalParts) != 2 {
		log.Fatalf("After splitting the suffix %s on `/` 2 parts were not found. The url should have /{webhook.id}/{webhook.token} at the end\n", urlParts[1])
	}
}

type WebhookMessage struct {
	Content string `json:"content"`
}

func postToWebhook(msg string) {
	reqMsg, err := json.Marshal(WebhookMessage{Content: msg})
	if err != nil {
		log.Printf("ERROR: Failed to marshal message to WebhookMessage.\nMessage: `%s`\n", msg)
		return
	}
	resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewReader(reqMsg))
	if err != nil {
		log.Printf("ERROR: Failed to post to webhook. %s", err.Error())
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Printf("ERROR: Status Code was not successful. got=%d", resp.StatusCode)
		return
	}
	log.Printf("Successfully Posted to Webhook!")
}

func handleSMTPConnection(conn net.Conn) {
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	writer.WriteString("220 smtp-server ESMTP ready\n")
	writer.Flush()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading: %v", err)
			return
		}

		log.Printf("RCV: %s", strings.TrimSpace(line))

		cmd := strings.ToUpper(strings.TrimSpace(strings.Replace(line, "\r\n", "", 1)))

		switch {
		case strings.HasPrefix(cmd, "HELO"), strings.HasPrefix(cmd, "EHLO"):
			writer.WriteString("250 smtp-server ready\n")
			writer.Flush()
		case strings.HasPrefix(cmd, "MAIL FROM:"):
			writer.WriteString("250 OK\n")
			writer.Flush()
		case strings.HasPrefix(cmd, "RCPT TO:"):
			writer.WriteString("250 OK\n")
			writer.Flush()
		case strings.HasPrefix(cmd, "DATA"):
			writer.WriteString("354 End data with <CR><LF>.<CR><LF>\n")
			writer.Flush()

			var msg strings.Builder
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					log.Printf("Error reading message: %v", err)
					return
				}

				if strings.TrimSpace(dataLine) == "." {
					break
				}
				msg.WriteString(dataLine)
			}

			msgToSend := msg.String()
			log.Printf("Received email:\n%s\n", msgToSend)
			postToWebhook(msgToSend)
			writer.WriteString("250 OK: queued\n")
			writer.Flush()
		case strings.HasPrefix(cmd, "QUIT"):
			writer.WriteString("221 Bye\n")
			writer.Flush()
			return
		default:
			writer.WriteString("502 Command not implemented\n")
			writer.Flush()
		}
	}
}

func main() {
	url = os.Getenv(webhookUrlEnvVariable)
	verifyUrlShape()

	listener, err := net.Listen("tcp", ":9123")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("SMTP server listening on port 9123...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		log.Println("\nShutting down...")
		listener.Close()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				// Just dont log on exit
				os.Exit(0)
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleSMTPConnection(conn)
	}
}
