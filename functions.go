package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getLoadAverage() string {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func toStr(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	case []byte:
		//return hex.EncodeToString(v)
		return string(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func getCPUTemperature() float64 {
	thermalZonePath := "/sys/class/thermal/thermal_zone0/temp"
	data, err := os.ReadFile(thermalZonePath)
	if err != nil {
		return 0
	}
	tempStr := strings.TrimSpace(string(data))
	tempMilli, err := strconv.Atoi(tempStr)
	if err != nil {
		return 0
	}
	return float64(tempMilli) / 1000.0
}

func msg2tlg(message string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Создаем данные для запроса
	reqBody := &SendMessageRequest{
		ChatID: chatID,
		Text:   message,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		echo(message)
	}
}

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
