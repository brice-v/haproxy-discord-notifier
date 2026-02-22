# Haproxy Discord Notifier

This program listens on port 9123 for smtp emails that are being sent from haproxy.
It then forwards those email notifications to Discord using webhooks


## Notes

Build `send-email` program for testing with `go build -o send-email cmd/send-email/main.go`
Build `haproxy-discord-notifier` with `go build`
Ctrl+C to kill notifier gracefully.
Include WEBHOOK_URL=... in running path (such as `WEBHOOK_URL=... ./haproxy-discord-notifier`) or just add it to your environment


### Sources

[This post](https://mko.re/blog/haproxy-webhook-alerts/)
[Haproxy docs](https://www.haproxy.com/documentation/haproxy-configuration-tutorials/alerts-and-monitoring/email-alerts/)
[Discord docs](https://docs.discord.com/developers/resources/webhook#execute-webhook)