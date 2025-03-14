package app

import (
	"context"
	"fmt"
)

func (bapp *BaseApp) startSystemd(ctx context.Context) error {
    command := fmt.Sprintf("sudo systemctl start %s", bapp.cfg.Node.ServiceName)
    err := runCommand(ctx, command)
    return err
}

func (bapp *BaseApp) stopSystemd(ctx context.Context) error {
    command := fmt.Sprintf("sudo systemctl stop %s", bapp.cfg.Node.ServiceName)
    err := runCommand(ctx, command)
    return err
}

func (bapp *BaseApp) nodeDBReset(ctx context.Context) error {
	command := fmt.Sprintf("mv %s/data/priv_validator_state.json %s", bapp.cfg.Node.NodeHome, bapp.cfg.Node.NodeHome)
	err := runCommand(ctx, command)

    target := fmt.Sprintf("%s/data/*", bapp.cfg.Node.NodeHome)
    if bapp.cfg.Node.IsWasm {
        target = fmt.Sprintf("%s %s/wasm", target, bapp.cfg.Node.NodeHome)
    }
    command = fmt.Sprintf("rm -rf %s", target)
    err = runCommand(ctx, command)

    command = fmt.Sprintf("mv %s/priv_validator_state.json %s/data", bapp.cfg.Node.NodeHome, bapp.cfg.Node.NodeHome)
    err = runCommand(ctx, command)

	return err
}

func (bapp *BaseApp) configureStateSync(ctx context.Context) error {
	command := fmt.Sprintf(`
SNAP_RPC="%s"

LATEST_HEIGHT=$(curl -s $SNAP_RPC/block | jq -r .result.block.header.height); \
BLOCK_HEIGHT=$((LATEST_HEIGHT - 2000)); \
TRUST_HASH=$(curl -s "$SNAP_RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

sed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\1true| ; \
s|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\1\"$SNAP_RPC,$SNAP_RPC\"| ; \
s|^(trust_height[[:space:]]+=[[:space:]]+).*$|\1$BLOCK_HEIGHT| ; \
s|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\1\"$TRUST_HASH\"|" %s/config/config.toml`,
		bapp.cfg.Node.RPC, bapp.cfg.Node.NodeHome)

	err := runCommand(ctx, command)
	return err
}

func (bapp *BaseApp) genSnapshot(ctx context.Context, snapshotFile string) error {
	target := "data"
	if bapp.cfg.Node.IsWasm {
        target += " wasm"
	}
    command := fmt.Sprintf("tar -C %s -cf - %s | lz4 - %s", bapp.cfg.Node.NodeHome, target, snapshotFile)

	err := runCommand(ctx, command)
    return err
}

func (bapp *BaseApp) removeSnapshot(ctx context.Context, snapshotFile string) error {
    command := fmt.Sprintf("rm %s", snapshotFile)

    err := runCommand(ctx, command)
    return err
}
