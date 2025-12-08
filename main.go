package main

import (
	"net/http"
)

func main() {
	const fileSep = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(fileSep)))

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	server.ListenAndServe()
}
