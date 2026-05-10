package main

type GameState struct {
	Played bool `json:"played"`
	Rating int  `json:"rating"`
}

type Game struct {
	ID       int    `json:"id"`
	Year     int    `json:"yr"`
	Title    string `json:"t"`
	Date     string `json:"d"`
	Platform string `json:"p"`
	Rel      string `json:"rel,omitempty"`
}

type StoredData struct {
	Games  []Game               `json:"games,omitempty"`
	State  map[string]GameState `json:"state,omitempty"`
	NextID int                  `json:"next_id,omitempty"`
}

type GameCreateRequest struct {
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Platform string `json:"platform"`
	Date     string `json:"date"`
	Rel      string `json:"rel,omitempty"`
}

type StateResponse struct {
	State map[string]GameState `json:"state"`
	Games []Game               `json:"games"`
}
