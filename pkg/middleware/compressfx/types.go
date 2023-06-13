package compressfx

import (
	"go.uber.org/fx"
	"io"
	"net/http"
)

type Params struct {
	fx.In
	Config Config
}

type Result struct {
	fx.Out
	Middleware Middleware
}

type Config struct {
	Level int `yaml:"level"`
}

type Middleware struct {
	level int
}

type compressResponseWriter struct {
	compressor io.Writer
	w          http.ResponseWriter
}
