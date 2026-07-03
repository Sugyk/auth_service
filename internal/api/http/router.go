package http_api

import (
	"context"
	"net/http"
)

type Handler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	server *http.Server
}

func NewRouter(handler Handler) *Router {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("POST /api/v1/auth/reg", handler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", handler.Login)
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
	return r.server.Shutdown(ctx)
}
