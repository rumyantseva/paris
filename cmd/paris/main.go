package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/rumyantseva/paris/internal/version"
)

func main() {
	logger := logrus.New().WithField("version", version.Version)

	logger.Infof(
		"The application [%v %v] is starting...",
		version.BuildTime,
		version.Commit,
	)

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("Business logic port is not set")
	}

	diagPort := os.Getenv("DIAG_PORT")
	if diagPort == "" {
		logger.Fatal("Diagnostics port is not ste")
	}

	r := mux.NewRouter()
	server := http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: r,
	}

	diagRouter := mux.NewRouter()
	diagRouter.HandleFunc("/health", func(
		w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	diagRouter.HandleFunc("/ready", func(
		w http.ResponseWriter, _ *http.Request) {
		logger.Info("Received ready request")
		time.Sleep(time.Minute)
		/// ...
		w.WriteHeader(http.StatusOK)
	})

	diag := http.Server{
		Addr:    net.JoinHostPort("", diagPort),
		Handler: diagRouter,
	}

	shutdown := make(chan error, 2)

	go func() {
		logger.Info("Business logic server is preparing...")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			shutdown <- err
		}
	}()

	go func() {
		logger.Info("Diagnostics server is preparing...")
		err := diag.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			shutdown <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		logger.Infof("Received `%v`.", x)

	case err := <-shutdown:
		logger.Infof("Received shutdown message: %v", err)
	}

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err := diag.Shutdown(timeout)
	if err != nil {
		logger.Error(err)
	}

	err = server.Shutdown(timeout)
	if err != nil {
		logger.Error(err)
	}
}
