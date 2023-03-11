package config

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/json"
	"github.com/gookit/config/v2/yaml"
	"github.com/rs/zerolog"
	"github.com/sicet7/go-compose-it/pkg/utils"
	"net"
	"os"
)

type LogConfiguration struct {
	Level string `mapstructure:"level" default:"1"`
	File  string `mapstructure:"file"`
}

type HttpConfiguration struct {
	Address          string           `mapstructure:"addr" default:"0.0.0.0:8080"`
	ShutdownWait     int              `mapstructure:"shutdown_wait" default:"60"`
	TlsConfiguration TlsConfiguration `mapstructure:"tls"`
	Net              NetConfiguration `mapstructure:"net"`
}

type TlsConfiguration struct {
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type NetConfiguration struct {
	TrustedProxies []string `mapstructure:"trustedProxies"`
}

type DatabaseConfiguration struct {
	Url string `mapstructure:"url" default:"sqlite:data.db"`
}

type Configuration struct {
	Log      LogConfiguration      `mapstructure:"log"`
	Database DatabaseConfiguration `mapstructure:"database"`
	Http     HttpConfiguration     `mapstructure:"http"`
}

var (
	conf   Configuration
	loaded = false
)

func (c Configuration) GetLogLevel() zerolog.Level {
	level, err := zerolog.ParseLevel(c.Log.Level)
	if err != nil {
		fmt.Printf("Failed to parse value passed into logLevel config: %v\n", err)
		os.Exit(1)
	}
	return level
}

func (c NetConfiguration) GetTrustedProxies() []*net.IPNet {
	var nets []*net.IPNet
	for _, cidrString := range c.TrustedProxies {
		_, cidr, err := net.ParseCIDR(cidrString)
		if err == nil {
			nets = append(nets, cidr)
		}
	}
	return nets
}

func LoadFromFile(configFile string) error {
	c := config.New("main")
	c.WithOptions(func(options *config.Options) {
		options.ParseEnv = true
		options.ParseDefault = true
		options.ParseKey = true
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

	// just to make sure error is triggered early
	conf.GetLogLevel()

	return nil
}

func Get() *Configuration {
	if !loaded {
		panic("cannot access configuration before initialization")
	}
	return &conf
}
