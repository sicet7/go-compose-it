package main

import (
	"github.com/sicet7/go-compose-it/src/config"
	"github.com/sicet7/go-compose-it/src/database"
	"github.com/sicet7/go-compose-it/src/database/entities"
	"github.com/sicet7/go-compose-it/src/database/repositories/user"
	"github.com/sicet7/go-compose-it/src/http"
	"github.com/sicet7/go-compose-it/src/http/handlers"
	"github.com/sicet7/go-compose-it/src/logger"
	"github.com/sicet7/go-compose-it/src/server"
	"go.uber.org/fx"
	"gorm.io/gorm"
	goHttp "net/http"
)

func main() {

	app := fx.New(
		fx.Provide(
			config.NewConfiguration,
			logger.NewLogger,
			database.NewConnection,
			user.NewRepository,
			server.NewHTTPServer,
			http.AsRoute(handlers.NewHealthHandler),
			fx.Annotate(
				http.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
		),
		fx.Invoke(func(db *gorm.DB) {
			err := db.AutoMigrate(entities.List()...)
			if err != nil {
				panic(err)
			}
		}),
		fx.Invoke(func(*goHttp.Server) {}),
	)

	app.Run()
}
