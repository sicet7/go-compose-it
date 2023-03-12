package create

import (
	"github.com/sicet7/go-compose-it/cmd/create/user"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(user.Command)
}

var Command = &cobra.Command{
	Use:  "create",
	Args: cobra.MinimumNArgs(1),
}
