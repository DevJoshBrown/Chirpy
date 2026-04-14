package main

import (
	"fmt"
	"net/http"
)

type Server struct {
}

func main() {
	fmt.Println("Running main.go")
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
