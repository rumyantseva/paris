package main

import (
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.Info("The aplication is starting...")

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

	diag := http.Server{
		Addr: net.JoinHostPort("", diagPort),
	}

	go func() {
		logger.Info("Business logic server is preparing...")
		server.ListenAndServe()
	}()

	logger.Info("Diagnostics server is preparing...")
	diag.ListenAndServe()
}
