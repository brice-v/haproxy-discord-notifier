package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

const (
	webhookUrlEnvVariable = "WEBHOOK_URL"
	expectedUrlPrefix     = "https://discord.com/api/webhooks/"
)

func verifyUrlShape(url string) {
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

			log.Printf("Received email:\n%s\n", msg.String())
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
	url := os.Getenv(webhookUrlEnvVariable)
	verifyUrlShape(url)

	listener, err := net.Listen("tcp", ":9123")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("SMTP server listening on port 9123...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleSMTPConnection(conn)
	}
}
