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

	go func() {
		logger.Info("Business logic server is preparing...")
		server.ListenAndServe()
	}()

	go func() {
		logger.Info("Diagnostics server is preparing...")
		diag.ListenAndServe()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	x := <-interrupt
	logger.Infof("Received `%v`. Application stopped.", x)

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	diag.Shutdown(timeout)
}
