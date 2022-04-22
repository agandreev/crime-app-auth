package server

import (
	"context"
	"net/http"
	"time"

	"github.com/agandreev/crime-app-auth/pkg/config"
)

type Server struct {
	s *http.Server
}

func NewServer(config config.ServerConfig, handler http.Handler) *Server {
	return &Server{
		s: &http.Server{
			Addr:              ":" + config.Port,
			Handler:           handler,
			ReadHeaderTimeout: time.Duration(config.ReadTimeout) * time.Second,
			WriteTimeout:      time.Duration(config.WriteTimeout) * time.Second,
			MaxHeaderBytes:    1 << config.MaxHeader,
		},
	}
}

func (s Server) Run() error {
	return s.s.ListenAndServe()
}

func (s Server) Stop(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
