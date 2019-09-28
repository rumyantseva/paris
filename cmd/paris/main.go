package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	diag := http.Server{
		Addr: ":8081",
	}

	go func() {
		server.ListenAndServe()
	}()

	diag.ListenAndServe()
}
