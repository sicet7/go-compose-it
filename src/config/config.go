package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/json"
	"github.com/gookit/config/v2/yaml"
	"github.com/rs/zerolog"
	"github.com/sicet7/go-compose-it/src/utils"
	"github.com/sicet7/go-compose-it/src/utils/env"
	"net"
	"os"
	"path/filepath"
)

type LogConfiguration struct {
	LevelString string `mapstructure:"level" default:"1"`
	File        string `mapstructure:"file"`
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

func (c LogConfiguration) Level() zerolog.Level {
	level, err := zerolog.ParseLevel(c.LevelString)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Failed to parse value passed into logLevel config: %v\n", err)))
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

func NewConfiguration() *Configuration {
	c := config.New("main")
	c.WithOptions(func(options *config.Options) {
		options.ParseEnv = true
		options.ParseDefault = true
		options.ParseKey = true
		options.EnableCache = true
	})
	c.AddDriver(yaml.Driver)
	c.AddDriver(json.Driver)

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	var files []string

	if utils.FileExists(exPath + "/config.yaml") {
		files = append(files, exPath+"/config.yaml")
	}

	envPath := env.ReadStringEnv("COMPOSE_IT_CONFIG_FILE", "")

	if envPath != "" && utils.FileExists(envPath) {
		files = append(files, envPath)
	}

	cliPath := flag.String("config", "", "--config=\"/path/to/config.yaml\"")

	flag.Parse()

	if *cliPath != "" && utils.FileExists(*cliPath) {
		files = append(files, *cliPath)
	}

	err = c.LoadFiles(files...)
	if err != nil {
		panic(err)
	}

	conf := Configuration{}
	err = c.Decode(&conf)

	if err != nil {
		panic(err)
	}

	conf.Log.Level()

	return &conf
}
