package main

import (
	"context"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/milkyway-labs/snapshot2s3/api"
	"github.com/milkyway-labs/snapshot2s3/app"
	"github.com/milkyway-labs/snapshot2s3/logger"
)

func main() {
	ctx := context.Background()

	f, err := os.ReadFile("config.toml")
	if err != nil {
		logger.Error(err)
		return
	}
	cfg := app.Config{}
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		logger.Error(err)
		return
	}

    apiServer := api.NewAPIServer(cfg.API.Port, cfg.Aws.Bucket, cfg.Aws.Region)
    bapp := app.NewBaseApp(&cfg, apiServer)

    go func() {
        err := apiServer.RunAPIServer()
        if err != nil {
            logger.Error(err)
            panic(err)
        }
    }()

    for {
        err = bapp.Run(ctx)
        if err != nil {
            logger.Error(err)
            panic(err)
        }

        time.Sleep(time.Duration(cfg.General.SnapshotInterval) * time.Hour)
    }
}
