package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/viniferr33/img-processor/pkg/logger"
	"go.uber.org/zap"
)

var (
	s *http.Server
)

type Config struct {
	Port    string
	Host    string
	Handler http.Handler
}

func Init(config Config) {
	s = &http.Server{
		Addr:    config.Host + ":" + config.Port,
		Handler: config.Handler,
	}
}

func Start() {
	if s == nil {
		logger.Fatal("server was not initialized")
		return
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, shutdownCancelCtx := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancelCtx()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal("graceful shutdown timed out... forcing exit")
			}
		}()

		err := s.Shutdown(shutdownCtx)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to close server", logger.FromError(err))
		}
		serverStopCtx()
	}()

	logger.Info("starting server", zap.String("addr", s.Addr))
	err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", logger.FromError(err))
	}
	logger.Info("shutting down server")

	<-serverCtx.Done()
}
