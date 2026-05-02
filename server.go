package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/RagnaCron/linko-starter/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	httpServer *http.Server
	store      store.Store
	cancel     context.CancelFunc
	logger     *slog.Logger
}

func newServer(store store.Store, port int, cancel context.CancelFunc, logger *slog.Logger) *server {
	mux := http.NewServeMux()

	s := &server{
		store:  store,
		cancel: cancel,
		logger: logger,
	}

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: metricsMiddleware(requestID()(requestLogger(logger)(mux))),
	}

	mux.HandleFunc("GET /", s.handlerIndex)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.Handle("POST /api/login", s.authMiddleware(http.HandlerFunc(s.handlerLogin)))
	mux.Handle("POST /api/shorten", s.authMiddleware(http.HandlerFunc(s.handlerShortenLink)))
	mux.Handle("GET /api/stats", s.authMiddleware(http.HandlerFunc(s.handlerStats)))
	mux.Handle("GET /api/urls", s.authMiddleware(http.HandlerFunc(s.handlerListURLs)))
	mux.HandleFunc("GET /{shortCode}", s.handlerRedirect)
	mux.HandleFunc("POST /admin/shutdown", s.handlerShutdown)

	return s
}

func (s *server) start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr(ln.Addr().Network(), ln.Addr().String())
	if err != nil {
		return err
	}

	s.logger.Debug(fmt.Sprintf("Linko is running on http://localhost:%d", addr.Port))

	if err := s.httpServer.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *server) shutdown(ctx context.Context) error {
	s.logger.Debug("Linko is shutting down")
	return s.httpServer.Shutdown(ctx)
}

func (s *server) handlerShutdown(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENV") == "production" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	go s.cancel()
}
