package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func ServeHTTP(logger *slog.Logger, srv *http.Server) error {
	doneCh := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Error("bad shutdown", "error", err)
		}
		close(doneCh)
	}()

	logger.Info("listening", "port", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-doneCh
	return nil
}
