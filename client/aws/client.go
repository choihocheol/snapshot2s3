package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func New(ctx context.Context, accessKeyID, secretKey, region, bucket string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

    client := &Client{
        client: s3.NewFromConfig(cfg),
        bucket: bucket,
    }
    
	return client, nil
}
