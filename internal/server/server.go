// Пакет для работы с сервером и регистрации обработчиков API.
package server

import (
	"GoExamGatewayAPI/internal/config"
	"GoExamGatewayAPI/internal/middleware"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"
)

const (
	news     = "news"
	comments = "comments"
	censor   = "censor"
)

const reqTime time.Duration = time.Second * 10

// Server - структура сервера.
type Server struct {
	srv   *http.Server
	mux   *http.ServeMux
	proxy map[string]string
	cl    *http.Client
}

// New - конструктор сервера.
func New(cfg *config.Config) *Server {
	p := make(map[string]string)
	p[news] = cfg.News
	p[comments] = cfg.Comments
	p[censor] = cfg.Censor

	m := http.NewServeMux()
	server := &Server{
		srv: &http.Server{
			Addr:         cfg.Address,
			Handler:      m,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		mux:   m,
		proxy: p,
		cl: &http.Client{
			Timeout: reqTime,
		},
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
	s.mux.HandleFunc("GET /news", News(s.proxy[news], s.cl))
	s.mux.HandleFunc("GET /news/id/{id}", NewsById(s.proxy[news], s.proxy[comments], s.cl))
	s.mux.HandleFunc("POST /comments/new", AddComment(s.proxy[comments], s.proxy[censor], s.cl))
}

// Shutdown останавливает сервер используя graceful shutdown.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop server: %s", err.Error())
	}
}
