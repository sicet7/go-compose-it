package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/sicet7/go-compose-it/src/utils"
	"go.uber.org/fx"
	"log"
	"net/http"
	"time"
)

type HttpServerConfig interface {
	HttpAddress() string
	HttpTlsCertFile() string
	HttpTlsKeyFile() string
}

func NewHTTPServer(
	conf HttpServerConfig,
	logger *zerolog.Logger,
	mux *http.ServeMux,
	lc fx.Lifecycle,
) *http.Server {

	certFile := conf.HttpTlsCertFile()
	keyFile := conf.HttpTlsKeyFile()

	tls := certFile != "" && keyFile != ""

	if tls && !utils.FileExists(certFile) {
		panic(errors.New(fmt.Sprintf("failed to locate cert file: \"%s\"\n", certFile)))
	}

	if tls && !utils.FileExists(keyFile) {
		panic(errors.New(fmt.Sprintf("failed to locate key file: \"%s\"\n", keyFile)))
	}

	server := http.Server{
		Addr:         conf.HttpAddress(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
		ErrorLog: log.New(
			logger.With().Str("type", "http-error").Logger(),
			"",
			log.Lmsgprefix|log.Llongfile,
		),
		Handler: mux,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				var err error
				if tls {
					err = server.ListenAndServeTLS(certFile, keyFile)
				} else {
					err = server.ListenAndServe()
				}

				if err != nil && err != http.ErrServerClosed {
					panic(errors.New(fmt.Sprintf("http server process failed: %v\n", err)))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
	return &server
}
