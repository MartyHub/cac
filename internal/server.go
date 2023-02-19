package internal

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const addr = "127.0.0.1:42000"

func Start() error {
	server := &http.Server{
		Addr:        addr,
		IdleTimeout: 30 * time.Second,
		ReadTimeout: 500 * time.Millisecond,
	}
	done := make(chan error, 1)

	go handleSignals(server, done)

	_ = server.ListenAndServe()

	err := <-done

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func handleSignals(server *http.Server, done chan<- error) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-sigs

	done <- server.Shutdown(context.Background())
}
