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
	"strconv"
	"sync"
)

//go:embed static
var staticFiles embed.FS

// ─── Domain ──────────────────────────────────────────────────────────────────

// GameState holds mutable user data for a single game entry.
type GameState struct {
	Played bool `json:"played"`
	Rating int  `json:"rating"`
}

// ─── Store ───────────────────────────────────────────────────────────────────

// store is a thread-safe, file-backed key/value map of game states.
// Keys are game IDs as decimal strings (matches JSON object key semantics).
type store struct {
	mu   sync.RWMutex
	path string
	data map[string]GameState
}

func newStore(path string) (*store, error) {
	s := &store{path: path, data: make(map[string]GameState)}

	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil // first run, empty state is fine
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &s.data); err != nil {
		return nil, err
	}
	slog.Info("store loaded", "path", path, "entries", len(s.data))
	return s, nil
}

// all returns a shallow copy of the state map, safe to use outside the lock.
func (s *store) all() map[string]GameState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]GameState, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out
}

// set updates or inserts the state for the given ID and flushes to disk.
func (s *store) set(id string, gs GameState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[id] = gs
	return s.flush()
}

// flush writes state to a temp file then atomically renames it.
func (s *store) flush() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// ─── Middleware ───────────────────────────────────────────────────────────────

// cors adds permissive CORS headers, useful during local development.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

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

// ─── Handlers ────────────────────────────────────────────────────────────────

func handleGetState(db *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, db.all())
	}
}

func handleUpdateGame(db *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if _, err := strconv.Atoi(idStr); err != nil {
			writeErr(w, http.StatusBadRequest, "id must be an integer")
			return
		}

		var gs GameState
		if err := json.NewDecoder(r.Body).Decode(&gs); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if gs.Rating < 0 || gs.Rating > 5 {
			writeErr(w, http.StatusBadRequest, "rating must be 0–5")
			return
		}

		if err := db.set(idStr, gs); err != nil {
			slog.Error("store.set", "id", idStr, "err", err)
			writeErr(w, http.StatusInternalServerError, "failed to persist state")
			return
		}

		writeJSON(w, http.StatusOK, gs)
	}
}

// ─── Main ────────────────────────────────────────────────────────────────────

func main() {
	db, err := newStore("state.json")
	if err != nil {
		slog.Error("store init failed", "err", err)
		os.Exit(1)
	}

	// Serve embedded static files, stripping the "static/" prefix.
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		slog.Error("embed.Sub", "err", err)
		os.Exit(1)
	}
	static := http.FileServer(http.FS(sub))

	mux := http.NewServeMux()

	// API routes (Go 1.22 method+pattern syntax).
	mux.HandleFunc("GET /api/state", handleGetState(db))
	mux.HandleFunc("PUT /api/games/{id}", handleUpdateGame(db))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Everything else goes to the static file server (SPA catch-all).
	mux.Handle("/", static)

	addr := ":8080"
	slog.Info("game log ready", "url", "http://localhost"+addr)
	if err := http.ListenAndServe(addr, cors(mux)); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
