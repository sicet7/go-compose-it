package cmd

import (
	"fmt"
	"github.com/sicet7/go-compose-it/pkg/config"
	"github.com/sicet7/go-compose-it/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "go-compose-it",
		Short: "A tool for hosting composer packages",
		Long:  "go-compose-it is a simple tool made for hosting composer packages",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "--config=\"/path/to/config.yaml\"")
}

func initConfig() {
	var err error
	if cfgFile == "" {
		cfgFile, err = utils.RequireStringEnv("COMPOSE_IT_CONFIG_FILE")
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
