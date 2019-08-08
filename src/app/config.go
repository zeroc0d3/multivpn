package app

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/viper"
)

// Config stores the application-wide configurations
var ConfigYml appConfig

type appConfig struct {
	// the path to the error message file. Defaults to "config/errors.yaml"
	ErrorFile string `mapstructure:"error_file"`

	// Binary OpenVPN Linux
	OpenvpnLinux string `mapstructure:"openvpn_linux"`
	// Binary OpenVPN Windows
	OpenvpnWindows string `mapstructure:"openvpn_windows"`
	// Path Keys OpenVPN (*.ovpn)
	PathKeys string `mapstructure:"path_keys"`

	// Configuration Keys
	FileName string `mapstructure:"file_name"`
	PathFile string `mapstructure:"path_file"`
	AuthFile string `mapstructure:"auth_file"`
}

func (config appConfig) Validate() error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.OpenvpnLinux, validation.Required),
		validation.Field(&config.OpenvpnWindows, validation.Required),
		validation.Field(&config.PathKeys, validation.Required),
	)
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
