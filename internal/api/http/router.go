package http_api

import (
	"context"
	"errors"
	"net/http"
)

type Handler interface{}

type Router struct {
	server *http.Server
}

func NewRouter(handler Handler) *Router {
	mux := http.NewServeMux()

	// Routes

	//

	server := &http.Server{
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
