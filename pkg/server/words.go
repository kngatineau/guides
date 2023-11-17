package server

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"

	sdk "go.flipt.io/flipt/sdk/go"
)

var words []string

func init() {
	w, err := os.ReadFile("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}

	words = strings.Split(string(w), "\n")

	for i := range words {
		j := rand.Intn(i + 1)
		words[i], words[j] = words[j], words[i]
	}

	words = words[:50000]
}

type Server struct {
	flipt sdk.SDK
}

func NewServer(flipt sdk.SDK) *Server {
	return &Server{flipt: flipt}
}

func (s *Server) ListWords(w http.ResponseWriter, r *http.Request) {
	words, err := getWords(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bubblesort(words)

	if err := json.NewEncoder(w).Encode(words); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func bubblesort(w []string) {
	for i := 0; i < len(w)-1; i++ {
		for j := 0; j < len(w)-i-1; j++ {
			if w[j] > w[j+1] {
				w[j], w[j+1] = w[j+1], w[j]
			}
		}
	}
}

func quicksort(w []string) {
	sort.Strings(w)
}

func getWords(ctx context.Context) ([]string, error) {
	return words, nil
}

func getUser(ctx context.Context) string {
	return ctx.Value("user").(string)
}
