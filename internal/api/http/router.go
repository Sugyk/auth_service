package http_api

import (
	"context"
	"errors"
	"net/http"
)

type Handler interface {
	Register(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	server *http.Server
}

func NewRouter(handler Handler) *Router {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("POST /reg", handler.Register)
	//

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &Router{
		server: server,
	}
}

func (r *Router) Start() error {
	return r.server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	var err error
	if err = r.server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
