package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func checkHost(host string) bool {
	_, err := net.DialTimeout("tcp", host+":80", 2*time.Second)
	if err != nil {
		return false
	}
	return true
}

func runCommand(cmd string) {
	command := exec.Command("bash", "-c", cmd)
	command.Run()
}

func echo(args ...interface{}) {

	fmt.Printf("%s", time.Now().Format("2006-01-02 15:04:05.000 "))

	if len(args) == 1 {
		fmt.Println(args...)
	} else {
		format := "%v"
		for i := 1; i < len(args); i++ {
			format += " %v"
		}
		format += "\n"
		fmt.Printf(format, args...)
	}
}

func getRemoteHost(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "remote ") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[1], nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("remote host not found!")
}

// check systemd unit
func systemd() {

	unit := "vpnSwitcher"
	path := "/etc/systemd/system/" + unit + ".service"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}

	text := `[Unit]
Description=` + unit + `
After=network.target

[Service]
Type=simple
ExecStart=`
	executablePath, _ := os.Executable()
	text += executablePath + "\n"
	text += "WorkingDirectory=" + filepath.Dir(executablePath) + "\n\n"
	text += `[Install]
WantedBy=multi-user.target
Alias=` + unit + `.service
`

	os.WriteFile(path, []byte(text), 0644)
	fmt.Println("Create unit [" + unit + "] in systemd. Run:\n\tsystemctl enable " + unit + "\n\tsystemctl start " + unit)
}
