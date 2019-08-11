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

package app

import (
	"fmt"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/viper"
)

// Config stores the application-wide configurations
var ConfigYml appConfig
var KeysYml appKeys

type appConfig struct {
	// the path to the error message file. Defaults to "config/errors.yaml"
	ErrorFile string `mapstructure:"error_file"`

	// Binary OpenVPN Linux
	OpenvpnLinux string `mapstructure:"openvpn_linux"`
	// Binary OpenVPN Windows
	OpenvpnWindows string `mapstructure:"openvpn_windows"`
}

type appKeys struct {
	// Configuration Keys
	FileName string `mapstructure:"file_name"`
	PathFile string `mapstructure:"path_file"`
	AuthFile string `mapstructure:"auth_file"`
}

func (config appConfig) Validate() error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.OpenvpnLinux, validation.Required),
		validation.Field(&config.OpenvpnWindows, validation.Required),
	)
}

func (key appKeys) ValidateKeys() error {
	//--> disable all validation
	return nil
	// return validation.ValidateStruct(&key,
	// 	validation.Field(&key.FileName, validation.Required),
	// 	validation.Field(&key.PathFile, validation.Required),
	// 	validation.Field(&key.AuthFile, validation.Required),
	// )
}

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
func LoadConfigYml(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.SetDefault("error_file", "./src/config/errors.yaml")

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}
	if err := v.Unmarshal(&ConfigYml); err != nil {
		return err
	}
	return ConfigYml.Validate()
}

// The keys file(s) should be named as keys.yaml.
func LoadKeysYml(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("keys")
	v.SetConfigType("yaml")

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}
	if err := v.Unmarshal(&KeysYml); err != nil {
		return err
	}

	// if arguments option is nil (not use option -> set to "default")
	var mapKeys = v.Get(os.Args[1])
	if mapKeys == "" {
		mapKeys = v.Get("default")
	}
	// fmt.Println(mapKeys)

	return KeysYml.ValidateKeys()
}
