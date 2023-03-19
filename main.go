package main

import (
	"flag"
	myApp "github.com/sicet7/go-compose-it/src/app"
	"github.com/sicet7/go-compose-it/src/config"
	"github.com/sicet7/go-compose-it/src/database"
	"github.com/sicet7/go-compose-it/src/database/entities"
	"github.com/sicet7/go-compose-it/src/database/repositories/user"
	"github.com/sicet7/go-compose-it/src/handlers"
	"github.com/sicet7/go-compose-it/src/http"
	"github.com/sicet7/go-compose-it/src/logger"
	"github.com/sicet7/go-compose-it/src/middlewares/accessLog"
	"github.com/sicet7/go-compose-it/src/middlewares/recovery"
	"github.com/sicet7/go-compose-it/src/server"
	"github.com/sicet7/go-compose-it/src/utils"
	"github.com/sicet7/go-compose-it/src/utils/env"
	"go.uber.org/fx"
	"gorm.io/gorm"
	goHttp "net/http"
	"os"
	"path/filepath"
)

func main() {

	app := fx.New(
		fx.Provide(
			http.AsRoute(handlers.NewHealthHandler),
			fx.Annotate(
				findConfigFiles,
				fx.ResultTags(`name:"config-files"`),
			),
			fx.Annotate(
				config.NewReader,
				fx.ParamTags(`name:"config-files"`),
			),
			fx.Annotate(
				myApp.NewConfiguration,
				fx.As(new(logger.LogConfig)),
				fx.As(new(database.ConnectionConfig)),
				fx.As(new(server.HttpServerConfig)),
			),
			logger.NewLogger,
			database.NewConnection,
			server.NewHTTPServer,
			user.NewRepository,
			middlewareInGroup(accessLog.NewAccessLogMiddleware, "global-middleware"),
			middlewareInGroup(recovery.NewRecoveryMiddleware, "global-middleware"),
			middlewareStackFromGroup("global-middleware-stack", "global-middleware"),
			fx.Annotate(
				http.NewServeMux,
				fx.ParamTags(`group:"routes"`, `name:"global-middleware-stack"`),
			),
		),
		fx.Invoke(autoMigrate),
		fx.Invoke(func(*goHttp.Server) {}),
	)

	app.Run()
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(entities.List()...)
	if err != nil {
		panic(err)
	}
}

func findConfigFiles() []string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	var files []string

	if utils.FileExists(exPath + "/config.yaml") {
		files = append(files, exPath+"/config.yaml")
	}

	envPath := env.ReadStringEnv("COMPOSE_IT_CONFIG_FILE", "")

	if envPath != "" && utils.FileExists(envPath) {
		files = append(files, envPath)
	}

	cliPath := flag.String("config", "", "--config=\"/path/to/config.yaml\"")

	flag.Parse()

	if *cliPath != "" && utils.FileExists(*cliPath) {
		files = append(files, *cliPath)
	}
	return files
}

func middlewareInGroup(f any, groupName string, paramTags ...string) any {
	if len(paramTags) > 0 {
		return fx.Annotate(
			f,
			fx.ParamTags(paramTags...),
			fx.As(new(http.Middleware)),
			fx.ResultTags(`group:"`+groupName+`"`),
		)
	}
	return fx.Annotate(
		f,
		fx.As(new(http.Middleware)),
		fx.ResultTags(`group:"`+groupName+`"`),
	)
}

func middlewareStackFromGroup(stackName string, groupName string) any {
	return fx.Annotate(
		func(globalMiddlewares []http.Middleware) http.Middleware {
			stack := http.NewMiddlewareStack()
			for _, middlewareContainer := range globalMiddlewares {
				stack.Add(middlewareContainer)
			}
			return stack
		},
		fx.ParamTags(`group:"`+groupName+`"`),
		fx.ResultTags(`name:"`+stackName+`"`),
	)
}
