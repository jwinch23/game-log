package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

func TestNewStore_NoFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := newStore(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.path != path {
		t.Fatalf("expected path %q, got %q", path, s.path)
	}
	if len(s.customGames) != 0 {
		t.Fatal("expected empty customGames")
	}
	if len(s.state) != 0 {
		t.Fatal("expected empty state")
	}
	if s.nextID != defaultNextID {
		t.Fatalf("expected nextID %d, got %d", defaultNextID, s.nextID)
	}
}

func TestNewStore_InvalidFormat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := newStore(path); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewStore_InvalidStateFormat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte(`"oops"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := newStore(path); err == nil {
		t.Fatal("expected error for invalid state file format")
	}
}

func TestNewStore_LegacyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	legacy := `{"123":{"played":true,"rating":5}}`
	if err := os.WriteFile(path, []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := newStore(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	got, ok := s.state[123]
	if !ok {
		t.Fatal("expected legacy entry to be loaded")
	}
	if got.Played != true || got.Rating != 5 {
		t.Fatalf("unexpected legacy state: %#v", got)
	}
}

func TestNewStore_StoredData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	payload := StoredData{
		Games:  []Game{{ID: 111, Year: 2021, Title: "A Game", Date: "2021", Platform: "ps"}},
		State:  map[string]GameState{"111": {Played: true, Rating: 4}},
		NextID: 2000,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := newStore(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.nextID != 2000 {
		t.Fatalf("expected nextID 2000, got %d", s.nextID)
	}
	if len(s.customGames) != 1 {
		t.Fatalf("expected one custom game, got %d", len(s.customGames))
	}
	if len(s.state) != 1 {
		t.Fatalf("expected one state entry, got %d", len(s.state))
	}
}

func TestStore_CreateGameAndSetState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := newStore(path)
	if err != nil {
		t.Fatal(err)
	}

	created, err := s.createGame(Game{Year: 2025, Title: "New Game", Date: "2025", Platform: "ps"})
	if err != nil {
		t.Fatalf("createGame error: %v", err)
	}
	if created.ID != defaultNextID {
		t.Fatalf("expected assigned ID %d, got %d", defaultNextID, created.ID)
	}

	provided := Game{ID: 42, Year: 2026, Title: "Manual ID", Date: "2026", Platform: "multi"}
	created2, err := s.createGame(provided)
	if err != nil {
		t.Fatalf("createGame error: %v", err)
	}
	if created2.ID != 42 {
		t.Fatalf("expected preserved ID 42, got %d", created2.ID)
	}

	if len(s.customGamesList()) != 2 {
		t.Fatalf("expected 2 custom games, got %d", len(s.customGamesList()))
	}

	expectedState := GameState{Played: true, Rating: 5}
	if err := s.setState(created.ID, expectedState); err != nil {
		t.Fatalf("setState error: %v", err)
	}

	mapped := s.stateMap()
	if got := mapped[strconv.Itoa(created.ID)]; !reflect.DeepEqual(got, expectedState) {
		t.Fatalf("unexpected state after setState: %#v", got)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var stored StoredData
	if err := json.Unmarshal(b, &stored); err != nil {
		t.Fatal(err)
	}
	if len(stored.Games) != 2 {
		t.Fatalf("expected 2 persisted games, got %d", len(stored.Games))
	}
	if stored.NextID != defaultNextID+1 {
		t.Fatalf("expected persisted next ID %d, got %d", defaultNextID+1, stored.NextID)
	}
}

func TestCreateGameRenameError(t *testing.T) {
	base := t.TempDir()
	preservedDir := filepath.Join(base, "state.json")
	if err := os.Mkdir(preservedDir, 0o755); err != nil {
		t.Fatal(err)
	}

	s := &store{path: preservedDir, nextID: defaultNextID, customGames: make(map[int]Game), state: make(map[int]GameState)}
	if _, err := s.createGame(Game{Year: 2025, Title: "Broken", Date: "2025", Platform: "ps"}); err == nil {
		t.Fatal("expected createGame to fail when rename target is a directory")
	}
}

func TestNewStore_StoredDataHighID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	payload := StoredData{
		Games: []Game{{ID: 2000, Year: 2024, Title: "High ID Game", Date: "2024", Platform: "ps"}},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := newStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.nextID != 2001 {
		t.Fatalf("expected nextID 2001, got %d", s.nextID)
	}
}

func TestNewStore_ReadFileError(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "state.json")
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := newStore(path); err == nil {
		t.Fatal("expected error when state path is a directory")
	}
}

func TestFlushLocked_MarshalError(t *testing.T) {
	old := jsonMarshalIndent
	defer func() { jsonMarshalIndent = old }()
	jsonMarshalIndent = func(v any, prefix, indent string) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}

	db, _ := newTemporaryStore(t)
	if err := db.setState(1, GameState{Played: true, Rating: 3}); err == nil {
		t.Fatal("expected error from failing marshaler")
	}
}
