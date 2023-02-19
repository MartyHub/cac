package internal

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const (
	addr = "127.0.0.1:42000"
)

func Start() error {
	pidFile, err := writePid()

	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:        addr,
		IdleTimeout: 30 * time.Second,
		ReadTimeout: 500 * time.Millisecond,
	}
	done := make(chan error, 1)

	go handleSignals(server, done)

	_ = server.ListenAndServe()
	_ = os.Remove(pidFile)

	err = <-done

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func GetPid() (int, error) {
	stateHome, err := GetStateHome()

	if err != nil {
		return 0, err
	}

	pidFile := filepath.Join(stateHome, "cac.pid")
	file, err := os.Open(pidFile)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(bytes))
}

func Stop() error {
	pid, err := GetPid()

	if err != nil {
		return err
	}

	p, err := os.FindProcess(pid)

	if err != nil {
		return err
	}

	defer func() {
		_ = p.Release()
	}()

	return p.Signal(syscall.SIGTERM)
}

func writePid() (string, error) {
	stateHome, err := GetStateHome()

	if err != nil {
		return "", err
	}

	pidFile := filepath.Join(stateHome, "cac.pid")
	file, err := os.OpenFile(pidFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(os.Getpid()))

	return pidFile, err
}

func handleSignals(server *http.Server, done chan<- error) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-sigs

	done <- server.Shutdown(context.Background())
}
