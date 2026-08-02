package http

import (
	"net/http"
	"time"
)

type ServerConfig struct {
	Addr              string
	ReadHeaderTimeout int
	ReadTimeout       int
	WriteTimeout      int
}

func NewServer(config *ServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    config.Addr,
		Handler: handler,

		ReadTimeout:       time.Duration(config.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(config.ReadHeaderTimeout) * time.Second,

		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
	}
}
