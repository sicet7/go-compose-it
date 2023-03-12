package cmd

import (
	"fmt"
	"github.com/sicet7/go-compose-it/cmd/create"
	"github.com/sicet7/go-compose-it/cmd/database"
	"github.com/sicet7/go-compose-it/cmd/serve"
	"github.com/sicet7/go-compose-it/pkg/config"
	"github.com/sicet7/go-compose-it/pkg/utils/env"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	cobra.OnInitialize(initConfig)

	command.AddCommand(serve.Command)
	command.AddCommand(create.Command)
	command.AddCommand(database.Command)

	command.PersistentFlags().StringVar(&cfgFile, "config", "", "--config=\"/path/to/config.yaml\"")
}

var (
	cfgFile string
	command = &cobra.Command{
		Use:   "go-compose-it",
		Short: "A tool for hosting composer packages",
		Long:  "go-compose-it is a simple tool made for hosting composer packages",
	}
)

func Execute() error {
	return command.Execute()
}

func initConfig() {
	var err error
	if cfgFile == "" {
		cfgFile, err = env.RequireStringEnv("COMPOSE_IT_CONFIG_FILE")
		if err != nil {
			fmt.Println("failed to find any configuration, use the --config option or the environment variable: " +
				"\"COMPOSE_IT_CONFIG_FILE\" to define the location of your configuration file.")
		}
	}

	err = config.LoadFromFile(cfgFile)
	if err != nil {
		fmt.Printf("failed to load config file: \"%s\"\n", cfgFile)
		os.Exit(1)
	}
}
