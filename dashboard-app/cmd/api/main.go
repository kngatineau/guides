package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"math/rand"
	"net/http"
	"os"

	"go.flipt.io/flipt/dashboard-app/ui"
	"go.flipt.io/flipt/rpc/flipt/evaluation"
	sdk "go.flipt.io/flipt/sdk/go"
	flipthttp "go.flipt.io/flipt/sdk/go/http"
)

var history = []struct {
	Date      string  `json:"date"`
	Sales     float64 `json:"Sales"`
	Profit    float64 `json:"Profit"`
	Customers int     `json:"Customers"`
}{
	{
		Date:      "2023-05-01",
		Sales:     900.73,
		Profit:    173,
		Customers: 73,
	},
	{
		Date:      "2023-05-02",
		Sales:     1000.74,
		Profit:    174.6,
		Customers: 74,
	},
	{
		Date:      "2023-05-03",
		Sales:     1100.93,
		Profit:    293.1,
		Customers: 293,
	},
	{
		Date:      "2023-05-04",
		Sales:     1200.9,
		Profit:    290.2,
		Customers: 29,
	},
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	flipt := sdk.New(flipthttp.NewTransport(fliptAddr()))

	dir, err := fs.Sub(ui.FS, "dist")
	if err != nil {
		slog.Error("embedding UI", "error", err)
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(http.FS(dir)))

	http.HandleFunc("/api/performance", func(w http.ResponseWriter, r *http.Request) {
		logger := slog.With(
			slog.String("namespace", "performance"),
			slog.String("flag", "showPerformanceHistory"),
		)

		// evaluate the showPerformanceHistory features flag
		result, err := flipt.Evaluation().Boolean(r.Context(), &evaluation.EvaluationRequest{
			NamespaceKey: "performance",
			FlagKey:      "showPerformanceHistory",
			EntityId:     fmt.Sprintf("%x", rand.Intn(1000)),
			Reference:    os.Getenv("FLIPT_CLIENT_REFERENCE"),
		})
		if err != nil {
			logger.Error("evaluating flag", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// if the flag is disabled we return that the endpoint cannot be found
		if !result.Enabled {
			logger.Debug("flag disabled")
			http.Error(w, "path not found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(&history); err != nil {
			logger.Error("parsing json", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	slog.Info("Listening", "addr", ":8081")

	http.ListenAndServe(":8081", nil)
}

func fliptAddr() string {
	if addr := os.Getenv("FLIPT_ADDRESS"); addr != "" {
		return addr
	}

	return "http://localhost:8080"
}
