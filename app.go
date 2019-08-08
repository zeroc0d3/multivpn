package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/dimiro1/banner"
	_ "github.com/dimiro1/banner/autoload"
	"github.com/mattn/go-colorable"
)

const MULTIVPN_LOG = "/var/log/multivpn/multivpn.log"

var loadKey string = "default.ovpn"
var cmdRunLinux string = "/usr/sbin/openvpn"
var cmdRunWindows string = "C:\\Program Files\\OpenVPN\\bin\\openvpn.exe"
var authFile string = "/opt/multivpn/keys/auth.txt"

func initLogo() {
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString("OpenVPN Selector {{ .AnsiColor.Green }}(Running){{ .AnsiColor.Default }}\n"))
}

func loadConfig() {
	// load option in yaml file
	//   --> load config/app.yaml    ; config binary & path
	//   --> load config/keys.yaml   ; config keys & auth (openvpn)
}

func runVPN() {
	var runKey string

	if runtime.GOOS == "windows" {
		runKey = fmt.Sprintf("%s --config %s --auth-user-pass %s", cmdRunWindows, loadKey, authFile)
	} else {
		runKey = fmt.Sprintf("%s --config %s --auth-user-pass %s", cmdRunLinux, loadKey, authFile)
	}
	cmd := exec.Command("%s", runKey)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("# Running VPN Access...                    ")
		fmt.Printf("[ TERMINATED ]\n")
	} else {
		fmt.Printf("# Running VPN Access...                    ")
		fmt.Printf("[   DONE   ]\n")
	}
}

func main() {
	initLogo()

	// load yaml file
	loadConfig()

	fmt.Printf("Running Key: %s \n", loadKey)
	fmt.Printf("Log file saved to: %s \n\n", MULTIVPN_LOG)

	// create the logger
	logger := logrus.New()
	// logger.Formatter = &logrus.JSONFormatter{}   # --> format to JSON

	fName, err := os.OpenFile(MULTIVPN_LOG, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer fName.Close()
	// multiwriter simultaneously
	logger.SetOutput(io.MultiWriter(os.Stdout, fName))

	runVPN()
}
