package app

import (
	"github.com/milkyway-labs/snapshot2s3/api"
	"github.com/milkyway-labs/snapshot2s3/client/aws"
)

type Config struct {
	General struct {
		SnapshotInterval int `toml:"snapshot_interval"`
	}
	API struct {
		Port string `toml:"port"`
	} `toml:"api"`
	Node struct {
		ServiceName string `toml:"service_name"`
		NodeHome    string `toml:"node_home"`
		RPC         string `toml:"rpc"`
		IsWasm      bool   `toml:"is_wasm"`
	} `toml:"node"`
	Aws struct {
		AccessKeyID     string `toml:"access_key_id"`
		SecretAccessKey string `toml:"secret_access_key"`
		Region          string `toml:"region"`
		Bucket          string `toml:"bucket"`
	} `toml:"aws"`
}

type BaseApp struct {
	cfg *Config
    apiServer *api.APIServer
    awsClient *aws.Client
}
