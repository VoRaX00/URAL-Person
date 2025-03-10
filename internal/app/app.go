package app

import (
	"context"
	"persons/internal/app/server"
)

type App struct {
	srv server.IServer
}

func NewApp(srv server.IServer) *App {
	return &App{
		srv: srv,
	}
}

func (a *App) MustStart() {
	if err := a.srv.Start(); err != nil {
		panic(err)
	}
}

func (a *App) MustStop(ctx context.Context) {
	if err := a.srv.Stop(ctx); err != nil {
		panic(err)
	}
}
