package config

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/json"
	"github.com/gookit/config/v2/yaml"
	"github.com/rs/zerolog"
	"go-compose-it/pkg/utils"
	"os"
)

type HttpConfiguration struct {
	Address          string           `mapstructure:"addr" default:"0.0.0.0:8080"`
	ShutdownWait     int              `mapstructure:"shutdown_wait" default:"60"`
	TlsConfiguration TlsConfiguration `mapstructure:"tls"`
}

type TlsConfiguration struct {
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type DatabaseConfiguration struct {
	Url string `mapstructure:"url" default:"sqlite:data.db"`
}

type Configuration struct {
	LogLevel zerolog.Level         `mapstructure:"logLevel"`
	LogFile  string                `mapstructure:"logFile"`
	Database DatabaseConfiguration `mapstructure:"database"`
	Http     HttpConfiguration     `mapstructure:"http"`
}

var (
	conf   Configuration
	loaded = false
)

func LoadFromFile(configFile string) error {
	c := config.New("main")
	c.WithOptions(func(options *config.Options) {
		options.ParseDefault = true
		options.ParseKey = true
		options.Readonly = true
		options.EnableCache = true
	})
	c.AddDriver(yaml.Driver)
	c.AddDriver(json.Driver)

	if !utils.FileExists(configFile) {
		fmt.Printf("Failed find configuration file: \"%s\"\n", configFile)
		os.Exit(1)
	}

	err := c.LoadFiles(configFile)
	if err != nil {
		return err
	}

	conf = Configuration{}
	err = c.Decode(&conf)
	if err != nil {
		return err
	}
	loaded = true
	return nil
}

func Get() *Configuration {
	if !loaded {
		panic("cannot access configuration before initialization")
	}
	return &conf
}
