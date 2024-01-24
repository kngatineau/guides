package main

import (
	"context"
	"log/slog"
	"net/http"

	"go.flipt.io/flipt/gitops-guide/pkg/server"
	_ "go.flipt.io/flipt/rpc/flipt"
	_ "go.flipt.io/flipt/rpc/flipt/evaluation"
	sdk "go.flipt.io/flipt/sdk/go"
	flipthttp "go.flipt.io/flipt/sdk/go/http"
)

func main() {
	t := flipthttp.NewTransport("http://localhost:8080")
	s := server.NewServer(sdk.New(t))

	http.HandleFunc("/words", threadUserContext(s.ListWords))
	slog.Info("Listening", "port", "8000")
	http.ListenAndServe(":8000", nil)
}

func threadUserContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		if user == "" {
			user = "default"
		}

		org := r.URL.Query().Get("org")
		if org == "" {
			org = "default"
		}

		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "org", org)

		next(w, r.WithContext(ctx))
	}
}
