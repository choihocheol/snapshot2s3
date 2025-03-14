package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/milkyway-labs/snapshot2s3/logger"
)

func NewAPIServer(port string, s3Bucket string, region string) *APIServer {
	return &APIServer{
		mux:           http.NewServeMux(),
		port:          fmt.Sprintf(":%s", port),
		s3url:         fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s3Bucket, region),
		SnapshotState: nil,
		AddrBookState: nil,
	}
}

func (apiServer *APIServer) NewState(fileName string, height int64, time time.Time) *State {
	return &State{
		RedirectURL: fmt.Sprintf("%s%s", apiServer.s3url, fileName),
		Height:      height,
		Time:        time,
	}
}

func (apiServer *APIServer) RunAPIServer() error {
	apiServer.mux.HandleFunc("/snapshot", apiServer.snapshotHandler)
	apiServer.mux.HandleFunc("/snapshot/status", apiServer.snapshotStatusHandler)
	apiServer.mux.HandleFunc("/addrbook", apiServer.addrbookHandler)
	apiServer.mux.HandleFunc("/addrbook/status", apiServer.addrbookStatusHandler)

	server := &http.Server{
		Addr:    apiServer.port,
		Handler: apiServer.mux,
	}

	msg := fmt.Sprintf("Starting server on %s", apiServer.port)
	logger.Info(msg)
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
