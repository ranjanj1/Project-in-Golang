package config

import (
	"flag"
	"io"
	"os"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	ServiceName string `mapstructure:"serviceName"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
}

const (
	configFileKey     = "configFile"
	defaultConfigFile = ""
	configFileUsage   = "this is config file path"
)

var (
	once         sync.Once
	cachedConfig AppConfig
)

type AppConfig struct {
	ServerConfig ServerConfig `mapstructure:"app"`
}

func (a *AppConfig) GetServerConfig() ServerConfig {
	return a.ServerConfig
}

func LoadConfig(reader io.Reader) (c AppConfig, err error) {

	keysToEnvironmentVariables := map[string]string{}

	err = loadConfig(reader, keysToEnvironmentVariables, &c)

	if err != nil {
		return c, err
	}

	return c, nil
}

func ProvideAppConfig() (c AppConfig, err error) {
	once.Do(func() {
		var configFile string
		flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
		flag.Parse()

		var configReader io.ReadCloser
		configReader, err = os.Open(configFile)
		defer configReader.Close() //nolint

		if err != nil {
			return
		}

		c, err = LoadConfig(configReader)
		if err != nil {
			return
		}

		cachedConfig = c
	})

	return cachedConfig, err
}

func loadConfig(reader io.Reader, bindings map[string]string, appConfig interface{}) error {

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := bind(bindings)
	if err != nil {
		return err
	}

	if err = viper.ReadConfig(reader); err != nil {
		return errors.Wrap(err, "Failed to load app config file")
	}

	if err = viper.Unmarshal(&appConfig); err != nil {
		return errors.Wrap(err, "Unable to parse app config file")
	}

	return nil
}

func bind(keysToEnvironmentVariables map[string]string) error {
	var bindErrors error
	for key, environmentVariable := range keysToEnvironmentVariables {
		if err := viper.BindEnv(key, environmentVariable); err != nil {
			bindErrors = multierror.Append(bindErrors, err)
		}
	}
	return bindErrors
}
