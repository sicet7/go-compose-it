package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/json"
	"github.com/gookit/config/v2/yaml"
)

func NewReader(files []string) *config.Config {
	c := config.New("main")
	c.WithOptions(func(options *config.Options) {
		//options.DecoderConfig.TagName = "config"
		options.ParseEnv = true
		options.ParseDefault = true
		options.ParseKey = true
		//options.EnableCache = true
	})
	c.AddDriver(yaml.Driver)
	c.AddDriver(json.Driver)

	err := c.LoadFiles(files...)
	if err != nil {
		panic(err)
	}
	return c
}
