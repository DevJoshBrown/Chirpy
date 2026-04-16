package main

//BUILD :  go build -o out && ./out
import (
	"fmt"
	"net/http"
)

type Server struct {
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	fmt.Println("Running main.go")
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlerReadiness)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	fileServerHandler := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", fileServerHandler))
	server.ListenAndServe()
}
