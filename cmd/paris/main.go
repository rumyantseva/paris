package main

import (
	"net/http"
)

func main() {
	server := http.Server{
		Addr: ":8080",
	}

	diag := http.Server{
		Addr: ":8081",
	}

	go func() {
		server.ListenAndServe()
	}()

	diag.ListenAndServe()
}
