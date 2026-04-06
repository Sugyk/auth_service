package http_api

import (
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
