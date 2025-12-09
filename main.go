package main

import (
	"net/http"
)

func main() {
	const fileSep = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(fileSep))))
	mux.HandleFunc("/healthz", healthz)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	server.ListenAndServe()
}
