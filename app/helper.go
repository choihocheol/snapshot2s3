package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/milkyway-labs/snapshot2s3/client/rpc"
	"github.com/milkyway-labs/snapshot2s3/logger"
)

func runCommand(ctx context.Context, command string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	msg := fmt.Sprintf("Running command: %s", command)
	logger.Info(msg)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Get the latest block height and whether the node is catching up
func (bapp *BaseApp) pollStatus(ctx context.Context, interval time.Duration) (int64, error) {
	rpcClient, err := rpc.New(bapp.cfg.Node.RPC)
	if err != nil {
		return 0, err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	status := &ctypes.ResultStatus{
		SyncInfo: ctypes.SyncInfo{
			CatchingUp: true,
		},
	}

	for status.SyncInfo.CatchingUp {
		select {
		case <-ticker.C:
			status, err = rpcClient.GetStatus(ctx)
			if err != nil {
				return 0, err
			}
		case <-ctx.Done():
			return 0, fmt.Errorf("Timeout during state-sync")
		}
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

func (bapp *BaseApp) uploadSnapshot(ctx context.Context, snapshotFileName string) error {
	length, oldest, err := bapp.awsClient.GetLengthAndOldestSnapshot(ctx)
	if err != nil {
		return err
	}

	// Only 2 snapshots are allowed to be stored in the bucket
	if length >= 2 {
		err = bapp.awsClient.DeleteFile(ctx, oldest)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Deleted the oldest snapshot: %s", oldest)
		logger.Info(msg)
	}

	err = bapp.awsClient.UploadFile(ctx, snapshotFileName, snapshotFileName)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Uploaded the snapshot: %s", snapshotFileName)
	logger.Info(msg)

	return nil
}

func (bapp *BaseApp) uploadAddrbook(ctx context.Context) error {
    // addrbook.json will overwrite even if it is existing
    addrbookPath := fmt.Sprintf("%s/config/%s", bapp.cfg.Node.NodeHome, "addrbook.json")
    err := bapp.awsClient.UploadFile(ctx, addrbookPath, "addrbook.json")
    if err != nil {
        return err
    }

    msg := fmt.Sprintf("Uploaded the addrbook: %s", addrbookPath)
    logger.Info(msg)

    return nil
}
