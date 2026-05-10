// game-log: personal game tracker HTTP server.
//
// Usage:
//   go run .
//   open http://localhost:8080
//
// State is persisted to ./state.json.
// Requires Go 1.22+ for enhanced routing patterns and PathValue.

package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
)

var (
	defaultStatePath = "state.json"
	listenAndServe   = http.ListenAndServe
	exitFunc         = os.Exit
	staticFS         fs.FS
)

//go:embed static
var staticFiles embed.FS

func main() {
	if err := run(defaultStatePath, ":8080", listenAndServe); err != nil {
		exitFunc(1)
	}
}

func run(statePath, addr string, serveFn func(string, http.Handler) error) error {
	db, err := newStore(statePath)
	if err != nil {
		slog.Error("store init failed", "err", err)
		return err
	}

	if staticFS == nil {
		staticFS = staticFiles
	}
	sub, _ := fs.Sub(staticFS, "static")
	if _, err := fs.Stat(sub, "."); err != nil {
		slog.Error("embed.Sub", "err", err)
		return err
	}
	static := http.FileServer(http.FS(sub))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/state", handleGetState(db))
	mux.HandleFunc("PUT /api/games/{id}", handleUpdateGame(db))
	mux.HandleFunc("POST /api/games", handleCreateGame(db))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("/", static)

	slog.Info("game log ready", "url", "http://localhost"+addr)
	if err := serveFn(addr, cors(mux)); err != nil {
		slog.Error("server error", "err", err)
		return err
	}
	return nil
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("writeJSON encode", "err", err)
	}
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
