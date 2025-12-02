package httputil

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server          *http.Server
	mux             *http.ServeMux
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

	return &Server{server, mux, config.ShutdownTimeout}
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
