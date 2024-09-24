package main

import (
	"os"
	"os/exec"
	"time"
)

func stat() {

	nowTime := time.Now().Unix()
	hostname, _ := os.Hostname()
	ifstat, _ := exec.Command("ifstat", "-a", "-b", "1", "1").Output()

	temp := getCPUTemperature()
	if temp != 0 && lastTimeMsg < (nowTime-sendPeriod) &&
		(temp < tempMin || temp > tempMax) {
		msg2tlg("[" + hostname + "] ⚠️  Temp: " + toStr(temp) + "\nLA: " + getLoadAverage() + "\n" + string(ifstat))
		lastTimeMsg = nowTime
	}

	// u can use another trigger for send info to telegram
	// msg2tlg("[" + hostname + "] warning!")

}
