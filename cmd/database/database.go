package database

import (
	"github.com/sicet7/go-compose-it/cmd/database/migrate"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(migrate.Command)
}

var Command = &cobra.Command{
	Use:   "database",
	Short: "Run commands against the database",
	Long:  "Run commands against the configured database",
}
