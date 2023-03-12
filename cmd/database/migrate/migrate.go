package migrate

import (
	"github.com/sicet7/go-compose-it/pkg/database"
	"github.com/sicet7/go-compose-it/pkg/models"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "migrate",
	Short: "Run migrations on the configured database",
	RunE:  command,
}

func command(cmd *cobra.Command, args []string) error {
	return database.RunMigrations(models.Get())
}
