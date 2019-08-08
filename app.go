package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/dimiro1/banner"
	_ "github.com/dimiro1/banner/autoload"
	"github.com/mattn/go-colorable"
)

const MULTIVPN_LOG = "/var/log/multivpn/multivpn.log"

var loadKey string = "ABC.ovpn"
var cmdRun string = "/usr/sbin/openvpn"

func initLogo() {
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString("OpenVPN Selector {{ .AnsiColor.Green }}(Running){{ .AnsiColor.Default }}\n"))
}

func loadConfig() {
	// load option in yaml file
}

func runVPN() {
	var runKey string

	runKey = fmt.Sprintf("%s --config %s", cmdRun, loadKey)
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
	logger.Formatter = &logrus.JSONFormatter{}

	fName, err := os.OpenFile(MULTIVPN_LOG, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer fName.Close()
	// multiwriter simultaneously
	logger.SetOutput(io.MultiWriter(os.Stdout, fName))

	runVPN()
}
