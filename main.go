package main

import (
	"github.com/sicet7/go-compose-it/src/config"
	"github.com/sicet7/go-compose-it/src/http"
	"github.com/sicet7/go-compose-it/src/http/handlers"
	"github.com/sicet7/go-compose-it/src/logger"
	"github.com/sicet7/go-compose-it/src/server"
	"go.uber.org/fx"
	goHttp "net/http"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.NewConfiguration,
			logger.NewLogger,
			server.NewHTTPServer,
			http.AsRoute(handlers.NewHealthHandler),
			fx.Annotate(
				http.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
		),
		fx.Invoke(func(*goHttp.Server) {}),
	)

	app.Run()
}
