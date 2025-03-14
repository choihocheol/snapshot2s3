package app

import (
	"context"
	"fmt"
	"time"

	"github.com/milkyway-labs/snapshot2s3/api"
	"github.com/milkyway-labs/snapshot2s3/client/aws"
	"github.com/milkyway-labs/snapshot2s3/logger"
)

func NewBaseApp(cfg *Config, apiServer *api.APIServer) *BaseApp {
	return &BaseApp{
		cfg:       cfg,
		apiServer: apiServer,
		awsClient: nil,
	}
}

func (bapp *BaseApp) Run(ctx context.Context) error {
	// Assume that all of the processes will be done in 3 hours
	// Consider the possibility that the state-sync failed during the process
	bappCtx, cancel := context.WithTimeout(ctx, 3*time.Hour)
	defer cancel()

	var err error

	// For interaction with the AWS S3
	bapp.awsClient, err = aws.New(
		bappCtx,
		bapp.cfg.Aws.AccessKeyID,
		bapp.cfg.Aws.SecretAccessKey,
		bapp.cfg.Aws.Region,
		bapp.cfg.Aws.Bucket,
	)
	if err != nil {
		return err
	}

	// Before start generating the snapshot, need to stop the node
	err = bapp.stopSystemd(bappCtx)
	if err != nil {
		return err
	}

	// 1. Reset node's db
	err = bapp.nodeDBReset(bappCtx)
	if err != nil {
		return err
	}

	// 2. Set state-sync configuration
	err = bapp.configureStateSync(bappCtx)
	if err != nil {
		return err
	}

	// 3. Start the systemd of node
	err = bapp.startSystemd(bappCtx)
	if err != nil {
		return err
	}

    logger.Info("Start state-sync")

	// 5. Wait for the state-sync to complete
	height, err := bapp.pollStatus(bappCtx, 1*time.Minute)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("State-sync completed successfully at height: %d", height)
	logger.Info(msg)

	current := time.Now().UTC()

	// 6. Stop the systemd of node
	err = bapp.stopSystemd(bappCtx)
	if err != nil {
		return err
	}

	snapshotFile := fmt.Sprintf("snapshot_%d.tar.lz4", height)

	// 7. Generate snapshot
	err = bapp.genSnapshot(bappCtx, snapshotFile)
	if err != nil {
		return err
	}

	// 8. Upload snapshot and addrbook to S3
	err = bapp.uploadSnapshot(bappCtx, snapshotFile)
	if err != nil {
		return err
	}
    err = bapp.uploadAddrbook(bappCtx)
    if err != nil {
        return err
    }

	// 9. Remove snapshot on local
	err = bapp.removeSnapshot(bappCtx, snapshotFile)
	if err != nil {
		return err
	}

	// 10. Change info for API server
	bapp.apiServer.SnapshotState = bapp.apiServer.NewState(snapshotFile, height, current)
    bapp.apiServer.AddrBookState = bapp.apiServer.NewState("addrbook.json", height, current)

	return nil
}
