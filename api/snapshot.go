package api

import (
	"encoding/json"
	"net/http"
)

func (apiServer *APIServer) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, apiServer.SnapshotState.RedirectURL, http.StatusFound)
}

func (apiServer *APIServer) snapshotStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
	json.NewEncoder(w).Encode(apiServer.SnapshotState)
}
