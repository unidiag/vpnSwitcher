package main

import (
	"os"
	"time"
)

const (
	botToken   = "XXXXXXXXXX:XXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXX" // telegramBot token
	chatID     = "XXXXXXXXX"                                      // recepient
	sendPeriod = 3600
	tempMin    = -10.0
	tempMax    = 70.0
)

type SendMessageRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

var vpnHost string
var lastTimeMsg int64

func main() {

	conf := "/etc/openvpn/client.conf" // default

	if len(os.Args) > 1 {
		conf = os.Args[1]
	}

	vpnHost, err := getRemoteHost(conf)
	if err != nil {
		echo("Unknown VPN host in " + conf)
		echo("Usage: " + os.Args[0] + " /etc/openvpn/client.conf")
		echo("\tor ")
		os.Exit(1)
	} else {
		echo("Run switcher VPN: " + vpnHost)
	}

	systemd()

	for {

		go check()
		go stat()

		time.Sleep(1 * time.Second)

	}
}
