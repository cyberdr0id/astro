package handler

import (
	"context"
	"net/http"

	"github.com/cyberdr0id/astro/internal/service"
	"github.com/gorilla/mux"
)

// Server a type that holds neccessary fields for running an HTTP server
type Server struct {
	server  *http.Server
	router  *mux.Router
	service service.APOD
}

// NewServer creates a new instance of Server.
func NewServer(service service.APOD) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		service: service,
	}

	s.router.HandleFunc("/apod", s.GetImage)
	s.router.HandleFunc("/entries", s.GetEntries)

	return s
}

// ServeHTTP dispatches the handler registered in the matched route.
func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

// Shutdown shuts down the server without interrupting any active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
