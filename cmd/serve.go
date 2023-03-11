package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"go-compose-it/pkg/config"
	myLogger "go-compose-it/pkg/logger"
	"go-compose-it/pkg/routes"
	"go-compose-it/pkg/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP Server",
	Long:  "Start the HTTP Server",
	Run: func(cmd *cobra.Command, args []string) {

		certFile := config.Get().Http.TlsConfiguration.CertFile
		keyFile := config.Get().Http.TlsConfiguration.KeyFile

		tls := certFile != "" && keyFile != ""

		if tls && !utils.FileExists(certFile) {
			fmt.Printf("failed to locate cert file: \"%s\"\n", certFile)
			os.Exit(1)
		}

		if tls && !utils.FileExists(keyFile) {
			fmt.Printf("failed to locate key file: \"%s\"\n", keyFile)
			os.Exit(1)
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

		// Start server in go routine
		go func() {
			var err error
			if tls {
				err = server.ListenAndServeTLS(certFile, keyFile)
			} else {
				err = server.ListenAndServe()
			}

			if err != nil && err != http.ErrServerClosed {
				fmt.Printf("failed to start http server: %v\n", err)
				os.Exit(1)
			}
		}()

		//Convert value int value from environment to a duration
		shutdownWait := time.Second * time.Duration(config.Get().Http.ShutdownWait)

		// Wait for signals to shut down gracefully
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		signal.Notify(c, syscall.SIGTERM)

		<-c

		ctx, cancel := context.WithTimeout(context.Background(), shutdownWait)
		defer cancel()

		err := server.Shutdown(ctx)

		if err != nil {
			fmt.Printf("http server shutdown threw errors: %v\n", err)
			os.Exit(1)
		}
	},
}
