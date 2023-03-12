package serve

import (
	"context"
	"errors"
	"fmt"
	"github.com/sicet7/go-compose-it/pkg/config"
	"github.com/sicet7/go-compose-it/pkg/database"
	myLogger "github.com/sicet7/go-compose-it/pkg/logger"
	"github.com/sicet7/go-compose-it/pkg/models"
	"github.com/sicet7/go-compose-it/pkg/routes"
	"github.com/sicet7/go-compose-it/pkg/utils"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	noMigration = false
	Command     = &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP Server",
		Long:  "Start the HTTP Server",
		RunE:  command,
	}
)

func init() {
	Command.Flags().BoolVarP(
		&noMigration,
		"no-migration",
		"",
		false,
		"--no-migration",
	)
}

func command(cmd *cobra.Command, args []string) error {
	certFile := config.Get().Http.TlsConfiguration.CertFile
	keyFile := config.Get().Http.TlsConfiguration.KeyFile

	tls := certFile != "" && keyFile != ""

	if tls && !utils.FileExists(certFile) {
		return errors.New(fmt.Sprintf("failed to locate cert file: \"%s\"\n", certFile))
	}

	if tls && !utils.FileExists(keyFile) {
		return errors.New(fmt.Sprintf("failed to locate key file: \"%s\"\n", keyFile))
	}

	server := http.Server{
		Addr:         config.Get().Http.Address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		ErrorLog: log.New(
			myLogger.Get("http-error"),
			"",
			log.Lmsgprefix|log.Llongfile,
		),
		Handler: routes.Mount(http.NewServeMux()),
	}

	if !noMigration {
		dbMigrationErr := database.RunMigrations(models.Get())

		if dbMigrationErr != nil {
			return errors.New(fmt.Sprintf("failed to run database migration: %v\n", dbMigrationErr))
		}
	}

	errCh := make(chan error, 1)

	// Start server in go routine
	go func() {
		var err error
		if tls {
			err = server.ListenAndServeTLS(certFile, keyFile)
		} else {
			err = server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			errCh <- errors.New(fmt.Sprintf("http server process failed: %v\n", err))
		}
	}()

	//Convert value int value from environment to a duration
	shutdownWait := time.Second * time.Duration(config.Get().Http.ShutdownWait)

	// Wait for signals to shut down gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	var err error

	select {
	case <-c:
		break
	case err = <-errCh:
	}

	if err != nil {
		return err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownWait)
		defer cancel()

		err = server.Shutdown(ctx)

		if err != nil {
			return errors.New(fmt.Sprintf("http server shutdown threw errors: %v\n", err))
		}
	}
	return nil
}
