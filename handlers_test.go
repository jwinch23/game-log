package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTemporaryStore(t *testing.T) (*store, string) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := newStore(path)
	if err != nil {
		t.Fatal(err)
	}
	return s, path
}

func TestHandleGetState(t *testing.T) {
	db, _ := newTemporaryStore(t)
	db.customGames[101] = Game{ID: 101, Year: 2023, Title: "Test Game", Date: "2023", Platform: "ps"}
	db.state[101] = GameState{Played: false, Rating: 2}

	req := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	resp := httptest.NewRecorder()
	handleGetState(db).ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var got StateResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.State["101"].Rating != 2 {
		t.Fatalf("expected rating 2, got %d", got.State["101"].Rating)
	}
	if len(got.Games) != 1 {
		t.Fatalf("expected 1 game, got %d", len(got.Games))
	}
}

func TestHandleUpdateGame(t *testing.T) {
	db, _ := newTemporaryStore(t)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/games/{id}", handleUpdateGame(db))

	req := httptest.NewRequest(http.MethodPut, "/api/games/123", strings.NewReader(`{"played":true,"rating":4}`))
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var got GameState
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Rating != 4 || !got.Played {
		t.Fatalf("unexpected state returned: %#v", got)
	}
}

func TestHandleUpdateGameInvalidInput(t *testing.T) {
	db, _ := newTemporaryStore(t)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/games/{id}", handleUpdateGame(db))

	cases := []struct {
		name string
		body string
		code int
	}{
		{"invalid-id", `{"played":true,"rating":3}`, http.StatusBadRequest},
		{"invalid-json", `{`, http.StatusBadRequest},
		{"bad-rating", `{"played":true,"rating":6}`, http.StatusBadRequest},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/games/100", strings.NewReader(c.body))
			if c.name == "invalid-id" {
				req = httptest.NewRequest(http.MethodPut, "/api/games/abc", strings.NewReader(c.body))
			}
			resp := httptest.NewRecorder()
			mux.ServeHTTP(resp, req)
			if resp.Code != c.code {
				t.Fatalf("expected %d, got %d", c.code, resp.Code)
			}
		})
	}
}

func TestHandleCreateGame(t *testing.T) {
	db, path := newTemporaryStore(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/games", handleCreateGame(db))

	payload := `{"title":"New Game","year":2025,"platform":"ps","date":"2025","rel":"2025-10-10"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(payload))
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.Code)
	}

	var got Game
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Title != "New Game" || got.Platform != "ps" {
		t.Fatalf("unexpected created game: %#v", got)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var stored StoredData
	if err := json.Unmarshal(b, &stored); err != nil {
		t.Fatal(err)
	}
	if len(stored.Games) != 1 {
		t.Fatalf("expected 1 persisted game, got %d", len(stored.Games))
	}
}

func TestHandleCreateGameInvalidRequest(t *testing.T) {
	db, _ := newTemporaryStore(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/games", handleCreateGame(db))

	cases := []struct {
		name string
		body string
	}{
		{"missing-title", `{"title":"","year":2025,"platform":"ps","date":"2025"}`},
		{"invalid-year", `{"title":"X","year":1800,"platform":"ps","date":"2025"}`},
		{"invalid-platform", `{"title":"X","year":2025,"platform":"pc","date":"2025"}`},
		{"missing-date", `{"title":"X","year":2025,"platform":"ps","date":""}`},
		{"bad-rel", `{"title":"X","year":2025,"platform":"ps","date":"2025","rel":"10-10-2025"}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(c.body))
			resp := httptest.NewRecorder()
			mux.ServeHTTP(resp, req)
			if resp.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", resp.Code)
			}
		})
	}
}

func TestHandleCreateGameInvalidJSON(t *testing.T) {
	db, _ := newTemporaryStore(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/games", handleCreateGame(db))

	req := httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(`{`))
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.Code)
	}
}

func TestHandleCreateGamePersistError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "state.json")
	db := &store{path: path, nextID: defaultNextID, customGames: make(map[int]Game), state: make(map[int]GameState)}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/games", handleCreateGame(db))

	req := httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(`{"title":"New Game","year":2025,"platform":"ps","date":"2025"}`))
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.Code)
	}
}

func TestHandleUpdateGamePersistError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "state.json")
	db := &store{path: path, nextID: defaultNextID, customGames: make(map[int]Game), state: make(map[int]GameState)}
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/games/{id}", handleUpdateGame(db))

	req := httptest.NewRequest(http.MethodPut, "/api/games/101", strings.NewReader(`{"played":true,"rating":3}`))
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.Code)
	}
}
