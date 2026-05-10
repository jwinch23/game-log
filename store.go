package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
)

const defaultNextID = 1000

var jsonMarshalIndent = json.MarshalIndent

type store struct {
	mu          sync.RWMutex
	path        string
	nextID      int
	customGames map[int]Game
	state       map[int]GameState
}

func newStore(path string) (*store, error) {
	s := &store{
		path:        path,
		nextID:      defaultNextID,
		customGames: make(map[int]Game),
		state:       make(map[int]GameState),
	}

	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	var data StoredData
	if err := json.Unmarshal(b, &data); err == nil && (len(data.Games) > 0 || len(data.State) > 0 || data.NextID > 0) {
		for _, g := range data.Games {
			s.customGames[g.ID] = g
			if g.ID >= s.nextID {
				s.nextID = g.ID + 1
			}
		}
		for k, v := range data.State {
			if id, err := strconv.Atoi(k); err == nil {
				s.state[id] = v
			}
		}
		if data.NextID > s.nextID {
			s.nextID = data.NextID
		}
		return s, nil
	}

	var legacy map[string]GameState
	if err := json.Unmarshal(b, &legacy); err == nil {
		for k, v := range legacy {
			if id, err := strconv.Atoi(k); err == nil {
				s.state[id] = v
			}
		}
		return s, nil
	}

	return nil, errors.New("invalid state file format")
}

func (s *store) stateMap() map[string]GameState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]GameState, len(s.state))
	for k, v := range s.state {
		out[strconv.Itoa(k)] = v
	}
	return out
}

func (s *store) customGamesList() []Game {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Game, 0, len(s.customGames))
	for _, g := range s.customGames {
		out = append(out, g)
	}
	return out
}

func (s *store) createGame(game Game) (Game, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if game.ID == 0 {
		game.ID = s.nextID
		s.nextID++
	}
	s.customGames[game.ID] = game
	if err := s.flushLocked(); err != nil {
		return Game{}, err
	}
	return game, nil
}

func (s *store) setState(id int, gs GameState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state[id] = gs
	return s.flushLocked()
}

func (s *store) flushLocked() error {
	rawState := make(map[string]GameState, len(s.state))
	for k, v := range s.state {
		rawState[strconv.Itoa(k)] = v
	}

	games := make([]Game, 0, len(s.customGames))
	for _, g := range s.customGames {
		games = append(games, g)
	}

	data := StoredData{
		Games:  games,
		State:  rawState,
		NextID: s.nextID,
	}

	b, err := jsonMarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
