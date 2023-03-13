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

	//bcryptHash := "$2y$10$AtjrsSrMYauotIz.Mnb8keJnh3z4g82HlVpP4uv20L6GtoacyJJ12"
	//argon2i := "$argon2i$v=19$m=65536,t=4,p=1$VEc5RkcvUWNyMEc4dDk4SA$DjfMPhGHjgT/dwbnEpzeyR+d2UZyVFWG/n4qw6T2zMM"
	//argon2id := "$argon2id$v=19$m=65536,t=4,p=1$MWRHVlNFVDFZcTdZQktDTg$xG3P+vF/fa/E/ovDaf8X3Hn4oE1/jak1bPQSJhOIt+Q"
	//pass := "test"
	//
	//h, err := password.Parse(bcryptHash)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//if h.VerifyPassword(pass) {
	//	fmt.Println("eyyyy")
	//} else {
	//	fmt.Println("fuck")
	//}
	//
	//h2, err2 := password.Parse(argon2i)
	//
	//if err2 != nil {
	//	panic(err2)
	//}
	//
	//if h2.VerifyPassword(pass) {
	//	fmt.Println("eyyyy")
	//} else {
	//	fmt.Println("fuck")
	//}
	//
	//h2, err2 = password.Parse(argon2id)
	//
	//if err2 != nil {
	//	panic(err2)
	//}
	//
	//if h2.VerifyPassword(pass) {
	//	fmt.Println("eyyyy")
	//} else {
	//	fmt.Println("fuck")
	//}

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
