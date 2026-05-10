package main

import (
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCorsOptionsAndPassThrough(t *testing.T) {
	h := cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp.Code)
	}
	if resp.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected CORS header to be set")
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d", resp.Code)
	}
}

func TestWriteJSONAndWriteErr(t *testing.T) {
	r := httptest.NewRecorder()
	writeJSON(r, http.StatusCreated, map[string]string{"hello": "world"})
	if r.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", r.Code)
	}
	if !strings.Contains(r.Body.String(), `"hello":"world"`) {
		t.Fatalf("unexpected body: %s", r.Body.String())
	}

	r = httptest.NewRecorder()
	writeErr(r, http.StatusBadRequest, "uh oh")
	if r.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", r.Code)
	}
	if !strings.Contains(r.Body.String(), `"error":"uh oh"`) {
		t.Fatalf("unexpected error body: %s", r.Body.String())
	}
}

func TestMainRunsWithInjectedListen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	oldListen := listenAndServe
	oldPath := defaultStatePath
	defer func() {
		listenAndServe = oldListen
		defaultStatePath = oldPath
	}()

	listenAndServe = func(addr string, handler http.Handler) error {
		return nil
	}
	defaultStatePath = path
	main()
}

func TestRunReturnsErrorWhenStateFileIsInvalid(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte("not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(path, ":0", func(addr string, handler http.Handler) error { return nil }); err == nil {
		t.Fatal("expected run to fail when state file is invalid")
	}
}

type failingFS struct{}

func (failingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("open failed")
}

func TestRunReturnsErrorWhenStaticFilesMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	oldStatic := staticFS
	defer func() { staticFS = oldStatic }()

	staticFS = failingFS{}
	if err := run(path, ":0", func(addr string, handler http.Handler) error { return nil }); err == nil {
		t.Fatal("expected run to fail when static files cannot be loaded")
	}
}

func TestMainCallsExitOnRunError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	oldListen := listenAndServe
	oldPath := defaultStatePath
	oldExit := exitFunc
	defer func() {
		listenAndServe = oldListen
		defaultStatePath = oldPath
		exitFunc = oldExit
	}()

	var exitCode int
	exitFunc = func(code int) {
		exitCode = code
	}
	listenAndServe = func(addr string, handler http.Handler) error {
		return errors.New("serve failed")
	}
	defaultStatePath = path
	main()
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
}

func TestWriteJSONEncodeFailure(t *testing.T) {
	r := httptest.NewRecorder()
	writeJSON(r, http.StatusOK, map[string]chan int{"bad": make(chan int)})
	if r.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", r.Code)
	}
}

func TestHealthzEndpoint(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	var capturedHandler http.Handler
	if err := run(path, ":0", func(_ string, h http.Handler) error {
		capturedHandler = h
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	resp := httptest.NewRecorder()
	capturedHandler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
