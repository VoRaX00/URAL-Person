package server

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Config struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type IServer interface {
	Start() error
	Stop(ctx context.Context) error
}

type Server struct {
	server *http.Server
}

func New(handler http.Handler, config Config) *Server {
	return &Server{
		server: &http.Server{
			Handler:      handler,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
