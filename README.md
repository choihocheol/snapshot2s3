# Snapshot2S3

This is a tool to generate a snapshot and upload snapshot and addrbook to S3 periodically.

## Architecture

![Architecture](asset/architecture.png)

## Features

- Upload snapshots to S3
    - Maintain only last 2 snapshots.
- Upload `addrbook.json` to S3
    - Maintain only last 1 `addrbook.json`.
- Efficiently generate snapshot
    - Only run a node during snapshot generating so that no need to run node continously.
 - Provide a API
    - `/snapshot`: Redirect to latest snapshot URL.
    - `/snapshot/status`: Get snapshot status related with `/snapshot` link.
        ```json
        {
          "redirect_url": "https://milkyway-mainnet-services.s3.us-west-1.amazonaws.com/snapshot_1565264.tar.lz4",
          "height": 1565264,
          "time": "2025-03-14T08:30:49.658703358Z"
        }
        ```
    - `/addrbook`: Redirect to latest `addrbook.json` URL.
    - `/addrbook/status`: Get `addrbook.json` status related with `/addrbook` link.
        ```json
        {
          "redirect_url": "https://milkyway-mainnet-services.s3.us-west-1.amazonaws.com/addrbook.json",
          "time": "2025-03-14T08:30:49.658703358Z"
        }
        ```

## Quick Guide

### Prerequisites

- `curl`, `jq`, `sed`, `tar`, `lz4` are required to be installed on the machine.
- This node will set up as systemd. (Just setup not need to run)
- AWS S3 bucket has to be configured public access.

### Usage

1. Build

```bash
go build
```

2. Configure `config.toml` file

```bash
cp config.toml.example config.toml
```

3. Run

```bash
./snapshot2s3 -config config.toml
```
