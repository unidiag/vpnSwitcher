package main

import "time"

func check() {

	cnt := 30 // every 30 sec
	add := "OK"

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
}
