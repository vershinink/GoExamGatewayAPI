// Пакет для работы с сервером и регистрации обработчиков API.
package server

import (
	"GoExamGatewayAPI/internal/config"
	"GoExamGatewayAPI/internal/middleware"
	"GoExamGatewayAPI/internal/server/api"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// Server - структура сервера.
type Server struct {
	srv *http.Server
	mux *http.ServeMux
}

// New - конструктор сервера.
func New(cfg *config.Config) *Server {
	m := http.NewServeMux()
	server := &Server{
		srv: &http.Server{
			Addr:         cfg.Address,
			Handler:      m,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		mux: m,
	}
	return server
}

// Start запускает HTTP сервер в отдельной горутине.
func (s *Server) Start() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			slog.Error("failed to start server")
		}
	}()
}

// Middleware инициализирует все обработчики middleware.
func (s *Server) Middleware() {
	wrappedMux := middleware.RealIP(middleware.RequestID(middleware.Logger(s.mux)))
	s.srv.Handler = wrappedMux
}

// API инициализирует все обработчики API.
func (s *Server) API() {
	s.mux.HandleFunc("GET /news/latest", api.Latest())
	s.mux.HandleFunc("GET /news/filter", api.Filter())
	s.mux.HandleFunc("GET /news/detailed/{id}", api.Detailed())
	s.mux.HandleFunc("POST /news/comment", api.AddComment())
}

// Shutdown останавливает сервер используя graceful shutdown.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop server: %s", err.Error())
	}
}
