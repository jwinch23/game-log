package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var validPlatforms = map[string]struct{}{
	"ps":       {},
	"nintendo": {},
	"multi":    {},
}

func handleGetState(db *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, StateResponse{
			State: db.stateMap(),
			Games: db.customGamesList(),
		})
	}
}

func handleUpdateGame(db *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
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

		if err := db.setState(id, gs); err != nil {
			slog.Error("store.setState", "id", id, "err", err)
			writeErr(w, http.StatusInternalServerError, "failed to persist state")
			return
		}

		writeJSON(w, http.StatusOK, gs)
	}
}

func handleCreateGame(db *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GameCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if err := validateCreateRequest(req); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}

		game := Game{
			Year:     req.Year,
			Title:    strings.TrimSpace(req.Title),
			Platform: req.Platform,
			Date:     strings.TrimSpace(req.Date),
			Rel:      strings.TrimSpace(req.Rel),
		}

		created, err := db.createGame(game)
		if err != nil {
			slog.Error("store.createGame", "err", err)
			writeErr(w, http.StatusInternalServerError, "failed to persist game")
			return
		}

		writeJSON(w, http.StatusCreated, created)
	}
}

func validateCreateRequest(req GameCreateRequest) error {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.Year < 1970 || req.Year > 2100 {
		return errors.New("year must be between 1970 and 2100")
	}
	if !isValidPlatform(req.Platform) {
		return errors.New("platform must be ps, nintendo, or multi")
	}
	if strings.TrimSpace(req.Date) == "" {
		return errors.New("display date is required")
	}
	if req.Rel != "" {
		if _, err := time.Parse("2006-01-02", req.Rel); err != nil {
			return errors.New("release date must be YYYY-MM-DD")
		}
	}
	return nil
}

func isValidPlatform(platform string) bool {
	_, ok := validPlatforms[platform]
	return ok
}
