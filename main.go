package main

import (
	"log/slog"
	"net/http"

	"go.flipt.io/flipt/gitops-guide/pkg/server"
)

func main() {
	s := &server.Server{}
	http.HandleFunc("/words", s.ListWords)
	slog.Info("Listening", "port", "8000")
	http.ListenAndServe(":8000", nil)
}
