package httputil

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apotourlyan/ludus-studii/pkg/httputil/middleware"
)

type Server struct {
	server          *http.Server
	mux             *http.ServeMux
	handler         http.Handler
	shutdownTimeout time.Duration
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

func NewServer(config *ServerConfig) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         config.Port,
		Handler:      mux,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	handler := middleware.CorrelationID(mux)

	return &Server{server, mux, handler, config.ShutdownTimeout}
}

func (s *Server) Handler() http.Handler {
	return s.handler
}

func (s *Server) AddMiddleware(m func(next http.Handler) http.Handler) {
	s.handler = m(s.handler)
}

func (s *Server) AddEndpoint(
	endpoint string, handler func(http.ResponseWriter, *http.Request),
) {
	s.mux.HandleFunc(endpoint, handler)
}

func (s *Server) Run() {
	go func() {
		log.Printf("Server listening on %s...", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	log.Println("Server stopped successfully")
}
