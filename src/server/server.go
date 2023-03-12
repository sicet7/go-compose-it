package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/sicet7/go-compose-it/src/config"
	"github.com/sicet7/go-compose-it/src/utils"
	"go.uber.org/fx"
	"log"
	"net/http"
	"time"
)

func NewHTTPServer(
	conf *config.Configuration,
	logger *zerolog.Logger,
	mux *http.ServeMux,
	lc fx.Lifecycle,
) *http.Server {

	certFile := conf.Http.TlsConfiguration.CertFile
	keyFile := conf.Http.TlsConfiguration.KeyFile

	tls := certFile != "" && keyFile != ""

	if tls && !utils.FileExists(certFile) {
		panic(errors.New(fmt.Sprintf("failed to locate cert file: \"%s\"\n", certFile)))
	}

	if tls && !utils.FileExists(keyFile) {
		panic(errors.New(fmt.Sprintf("failed to locate key file: \"%s\"\n", keyFile)))
	}

	server := http.Server{
		Addr:         conf.Http.Address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		ErrorLog: log.New(
			logger.With().Str("logger", "http-error").Logger(),
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
