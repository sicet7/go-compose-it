package app

import (
	"errors"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/rs/zerolog"
	"net"
	"os"
)

type LogConfiguration struct {
	LevelString string `mapstructure:"level" default:"1"`
	FilePath    string `mapstructure:"file"`
}

type HttpConfiguration struct {
	AddressString      string   `mapstructure:"addr" default:"0.0.0.0:8080"`
	CompLevelString    int      `mapstructure:"compression_level" default:"9"`
	TrustedProxiesList []string `mapstructure:"trustedProxies"`
	TlsCertFileString  string   `mapstructure:"tls_cert_file"`
	TlsKeyFileString   string   `mapstructure:"tls_key_file"`
}

type NetConfiguration struct {
	TrustedProxies []string `mapstructure:"trustedProxies"`
}

type Configuration struct {
	LogLevelString       string `mapstructure:"log.level" default:"1"`
	LogFilePath          string `mapstructure:"log.file"`
	DbUrl                string `mapstructure:"database.url" default:"sqlite:data.db"`
	HttpAddr             string `mapstructure:"http.addr" default:"0.0.0.0:8080"`
	HttpTlsCertFilePath  string `mapstructure:"http.tls.cert_file"`
	HttpTlsKeyFilePath   string `mapstructure:"http.tls.key_file"`
	HttpCompressionLevel int    `mapstructure:"http.compress.level" default:"9"`
}

func (c Configuration) CompressionLevel() int {
	return c.HttpCompressionLevel
}

func (c Configuration) HttpAddress() string {
	return c.HttpAddr
}

func (c Configuration) HttpTlsCertFile() string {
	return c.HttpTlsCertFilePath
}

func (c Configuration) HttpTlsKeyFile() string {
	return c.HttpTlsKeyFilePath
}

func (c Configuration) DatabaseUrl() string {
	return c.DbUrl
}

func (c Configuration) LogLevel() zerolog.Level {
	level, err := zerolog.ParseLevel(c.LogLevelString)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Failed to parse value passed into logLevel config: %v\n", err)))
	}
	return level
}

func (c Configuration) LogFile() *os.File {
	if c.LogFilePath != "" {
		openLogFile, err := os.OpenFile(c.LogFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		return openLogFile
	}
	return os.Stdout
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

func NewConfiguration(c *config.Config) *Configuration {
	conf := Configuration{}
	err := c.Decode(&conf)

	if err != nil {
		panic(err)
	}

	conf.LogLevel()

	return &conf
}
