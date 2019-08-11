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
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Sirupsen/logrus"
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

var loadKey string
var authFile string
var runMultivpn string
var str_name_file string
var str_path_file string
var str_auth_file string

func loadConfig() {
	// load configuration in environment variables:
	env := os.Getenv("ENV_MULTIVPN")
	if "development" == env {
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

	// if arguments[1] doesn't exist then use "default" key
	optArgs := os.Args[1]
	// fmt.Println(optArgs)

	if optArgs != "" {
		// load keys openvpn
		//   --> load config/keys.yaml   ; config keys & auth (openvpn)
		if err := app.LoadKeysYml(MULTIVPN_PATH_CONFIG); err != nil {
			panic(fmt.Errorf("Invalid application keys: %s", err))
		}

		str_name_file = optArgs + ".name_file"
		str_path_file = optArgs + ".path_file"
		str_auth_file = optArgs + ".auth_file"

		loadKey = app.KeysYml.PathFile + app.KeysYml.FileName
		authFile = app.KeysYml.AuthFile
	} else {
		// default
		loadKey = MULTIVPN_PATH_KEYS + MULTIVPN_DEFAULT_KEYS
		authFile = MULTIVPN_PATH_KEYS + MULTIVPN_DEFAULT_AUTH
	}

	if runtime.GOOS == "windows" {
		runMultivpn = fmt.Sprintf("%s --config %s --auth-user-pass %s", OPENVPN_BIN_WINDOWS, loadKey, authFile)
	} else {
		runMultivpn = fmt.Sprintf("%s --config %s --auth-user-pass %s", OPENVPN_BIN_LINUX, loadKey, authFile)
	}
}

func runVPN() {
	//runMultivpn := []rune(runMultivpn)
	//fmt.Println(string(runMultivpn[0:6]))
	_, err := exec.Command(runMultivpn).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("# Running VPN Access...                    ")
		fmt.Printf("[ TERMINATED ]\n")
	} else {
		fmt.Printf("# Running VPN Access...                    ")
		fmt.Printf("[   DONE   ]\n")
	}
}

func multivpnExecute() {
	// load yaml file
	loadConfig()
	fmt.Printf("OpenVPN Key : %s \n", loadKey)
	fmt.Printf("Auth File   : %s \n", authFile)
	fmt.Printf("Running     : %s \n", runMultivpn)
	fmt.Printf("Log File    : %s \n", MULTIVPN_LOG)
	fmt.Println("----------------------------------------------------------------------------")
	// running
	// runVPN()
}
