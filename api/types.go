package api

import (
	"net/http"
	"time"
)

type APIServer struct {
	mux           *http.ServeMux
	port          string
	s3url         string
	SnapshotState *State
	AddrBookState *State
}

// This struct using json response too.
type State struct {
	RedirectURL string    `json:"redirect_url"`
	Height      int64     `json:"height"`
	Time        time.Time `json:"time"`
}
