package main

import (
	"os"
	"time"
)

func main() {

	conf := "/etc/openvpn/client.conf" // default

	if len(os.Args) > 1 {
		conf = os.Args[1]
	}

	cnt := 30 // every 30 sec
	add := ""
	vpnHost, err := getRemoteHost(conf)
	if err != nil {
		echo("Unknown VPN host in " + conf)
		echo("Usage: " + os.Args[0] + " /etc/openvpn/client.conf")
		os.Exit(1)
	} else {
		echo("Run switcher VPN: " + vpnHost)
	}

	systemd()

	for {

		add = "OK"

		if checkHost("10.8.0.1") {

			if !checkHost("ya.ru") {
				if cnt == 0 {
					runCommand("systemctl stop openvpn")
					echo("No internet! Stopped openvpn...")
					time.Sleep(5 * time.Second)
					cnt = 30
				} else {
					cnt--
				}
				add = "NOK"
			}

			echo("VPN online. Internet: " + add)

		} else if checkHost(vpnHost) {

			if cnt == 0 {
				runCommand("systemctl start openvpn")
				echo("Try start openvpn...")
				time.Sleep(5 * time.Second)
				cnt = 30
			} else {
				cnt--
			}

			echo("VPN offline. Internet: OK")

		} else if !checkHost("ya.ru") {

			echo("VPN offline. Internet: NOK")

		}

		time.Sleep(1 * time.Second)

	}
}
