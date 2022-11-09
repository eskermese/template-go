package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eskermese/template-go/pkg/logger"

	"github.com/eskermese/template-go/internal/config"
)

type Server struct {
	CertFile, KeyFile *string
	httpServer        *http.Server
	masterCtx         context.Context
	idleConnsClosed   chan struct{}
	logger            logger.Logger
}

func NewServer(ctx context.Context, cfg *config.Config, handler http.Handler, logger logger.Logger) *Server {
	return &Server{
		masterCtx:       ctx,
		idleConnsClosed: make(chan struct{}),
		httpServer: &http.Server{
			Addr:           fmt.Sprintf(":%d", cfg.ClientHTTP.Port),
			Handler:        handler,
			ReadTimeout:    cfg.ClientHTTP.ReadTimeout,
			WriteTimeout:   cfg.ClientHTTP.WriteTimeout,
			MaxHeaderBytes: cfg.ClientHTTP.MaxHeaderMegabytes << 20,
		},
		logger: logger,
	}
}

func (s *Server) IsInsecure() bool {
	return s.CertFile == nil && s.KeyFile == nil
}

func (s *Server) Run() error {
	s.logger.Info(fmt.Sprintf(`serving ClientHTTP on "%s"`, s.httpServer.Addr))

	go s.gracefulShutdown(s.httpServer)

	var err error
	if s.IsInsecure() {
		err = s.httpServer.ListenAndServe()
	}

	if !s.IsInsecure() {
		err = s.httpServer.ListenAndServeTLS(*s.CertFile, *s.KeyFile)
	}

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	s.Wait()

	return s.httpServer.ListenAndServe()
}

func (s *Server) gracefulShutdown(httpSrv *http.Server) {
	defer close(s.idleConnsClosed)
	<-s.masterCtx.Done()

	s.logger.Info("shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), s.httpServer.ReadTimeout+s.httpServer.WriteTimeout)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		s.logger.Error("shutting down HTTP server", logger.Error(err))
	}
}

func (s *Server) Wait() {
	<-s.idleConnsClosed
	s.logger.Info("HTTP server has processed all idle connections")
}
