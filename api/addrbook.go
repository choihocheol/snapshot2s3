package api

import (
	"encoding/json"
	"net/http"
)

func (apiServer *APIServer) addrbookHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, apiServer.AddrBookState.RedirectURL, http.StatusFound)
}

func (apiServer *APIServer) addrbookStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
	json.NewEncoder(w).Encode(apiServer.AddrBookState)
}
