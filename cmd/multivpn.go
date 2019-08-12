/*
Copyright © 2019 ZeroLabs <zeroc0d3.team@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
	"github.com/zeroc0d3/multivpn/src/app"
	"github.com/zeroc0d3/multivpn/src/errors"
)

const MULTIVPN_LOG = "/var/log/multivpn/multivpn.log"

var MULTIVPN_PATH_CONFIG string = "/opt/multivpn/config/"
var MULTIVPN_PATH_KEYS string = "/opt/multivpn/keys/"
var MULTIVPN_DEFAULT_KEYS string = "default.ovpn"
var MULTIVPN_DEFAULT_AUTH string = "auth.txt"

var OPENVPN_BIN_LINUX string
var OPENVPN_BIN_WINDOWS string

var environment string
var loadKey string
var authFile string
var err error

func initLogo() {
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString("MultiVPN CLI {{ .AnsiColor.Green }}(Running){{ .AnsiColor.Default }} ...\n\n"))
}

func loadConfig() {
	// load configuration in environment variables:
	environment := os.Getenv("ENV_MULTIVPN")
	if "development" == environment {
		MULTIVPN_PATH_CONFIG = "./src/config/"
		MULTIVPN_PATH_KEYS = "./keys/"
	}

	// load configuration yaml file
	//   --> load config/app.yaml    ; config binary & path
	if err := app.LoadConfigYml(MULTIVPN_PATH_CONFIG); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}
	// load error messages
	if err := errors.LoadMessages(app.ConfigYml.ErrorFile); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	// create the logger
	logger := logrus.New()
	// logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(os.Stdout)

	file, err := os.OpenFile(MULTIVPN_LOG, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	logger.SetOutput(file)

	// locate binary openvpn
	if app.ConfigYml.OpenvpnLinux != "" {
		OPENVPN_BIN_LINUX = app.ConfigYml.OpenvpnLinux
	}
	if app.ConfigYml.OpenvpnWindows != "" {
		OPENVPN_BIN_WINDOWS = app.ConfigYml.OpenvpnWindows
	}

	if len(os.Args[2:]) < 1 {
		fmt.Println("Please provide at least one argument, to see available argument just type -h argument")
		os.Exit(1)
	}

	// load keys openvpn
	//   --> load config/keys.yaml   ; config keys & auth (openvpn)
	if err := app.LoadKeysYml(MULTIVPN_PATH_CONFIG); err != nil {
		panic(fmt.Errorf("Invalid application keys: %s", err))
	}

	if (app.KeysYml.PathFile != "") && (app.KeysYml.FileName != "") {
		loadKey = app.KeysYml.PathFile + app.KeysYml.FileName
	} else {
		// loadKey = MULTIVPN_PATH_KEYS + MULTIVPN_DEFAULT_KEYS
		fmt.Println(">> Can't use your openvpn (*.ovpn) key ...")
	}

	if app.KeysYml.AuthFile != "" {
		authFile = app.KeysYml.AuthFile
	} else {
		// authFile = MULTIVPN_PATH_KEYS + MULTIVPN_DEFAULT_AUTH
		fmt.Println(">> Can't use your auth configuration file ...\n")
	}
}

func runVPN() {
	var runBinary string
	if runtime.GOOS == "windows" {
		// runMultivpn = fmt.Sprintf("%s --config %s --auth-user-pass %s", OPENVPN_BIN_WINDOWS, loadKey, authFile)
		runBinary = OPENVPN_BIN_WINDOWS
	} else {
		// runMultivpn = fmt.Sprintf("%s --config %s --auth-user-pass %s", OPENVPN_BIN_LINUX, loadKey, authFile)
		runBinary = OPENVPN_BIN_LINUX
	}
	args := []string{"--config", loadKey, "--auth-user-pass", authFile}
	if err := exec.Command(runBinary, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Printf("# Running VPN Access...                    ")
		fmt.Printf("[ TERMINATED ]\n")
		os.Exit(1)
	} else {
		fmt.Printf("# Running VPN Access...                    ")
		fmt.Printf("[   DONE   ]\n")
		os.Exit(0)
	}
}

func multivpnExecute() {
	var env string = "development"
	initLogo()
	// load yaml file
	loadConfig()

	environment := os.Getenv("ENV_MULTIVPN")
	// fmt.Printf("env: %s", os.Getenv("ENV_MULTIVPN"))
	if "development" != environment {
		env = "production"
	}
	if loadKey != "" && authFile != "" {
		fmt.Println("----------------------------------------------------------------------------")
		fmt.Printf("Environment : %s \n", env)
		fmt.Printf("OpenVPN Key : %s \n", loadKey)
		fmt.Printf("Auth File   : %s \n", authFile)
		fmt.Printf("Log File    : %s \n", MULTIVPN_LOG)
		fmt.Println("----------------------------------------------------------------------------")
	}
	// running
	runVPN()
}
