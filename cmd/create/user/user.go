package user

import (
	"errors"
	"fmt"
	"github.com/sicet7/go-compose-it/pkg/database"
	"github.com/sicet7/go-compose-it/pkg/models"
	termUtils "github.com/sicet7/go-compose-it/pkg/utils/term"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"gorm.io/gorm"
	"os"
	"syscall"
)

var Command = &cobra.Command{
	Use:  "user",
	RunE: command,
}

func command(cmd *cobra.Command, args []string) error {

	if !term.IsTerminal(syscall.Stdin) {
		return errors.New("this command does not support piped data")
	}

	username := termUtils.StringPrompt("Enter username: ", os.Stdin, os.Stdout)
	password := termUtils.PasswordPrompt("Enter password: ", syscall.Stdin, os.Stdout)

	_, err := models.FindUserByUsername(username)

	if err == nil {
		return errors.New(fmt.Sprintf("a user with the username \"%s\" already exists", username))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("user lookout failed with error: %v", err))
	}

	user, createErr := models.NewUser(username, password)

	if createErr != nil {
		return errors.New(fmt.Sprintf("user creation failed with error: %v", createErr))
	}

	result := database.Conn().Create(&user)

	if result.Error != nil {
		return errors.New(fmt.Sprintf("user insert failed with error: %v", result.Error))
	}

	fmt.Println("user was successfully created")
	return nil
}
