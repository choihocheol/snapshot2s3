package aws

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Client) UploadFile(ctx context.Context, filePath string, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.bucket,
		Body:   file,
		Key:    &fileName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteFile(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: &c.bucket,
			Key:    &key,
		})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetLengthAndOldestSnapshot(ctx context.Context) (int, string, error) {
	resp, err := c.client.ListObjectsV2(
		ctx,
		&s3.ListObjectsV2Input{
			Bucket: &c.bucket,
		})
	if err != nil {
		return 0, "", err
	}

	snapshotCnt := 0

	// For setting the oldest object to the first object in the bucket
    tempKey := ""
	t := time.Now().UTC().Add(time.Hour * 24)
	oldestObject := types.Object{
		Key:          &tempKey,
		LastModified: &t,
	}
	for _, object := range resp.Contents {
		if *object.Key != "addrbook.json" {
			snapshotCnt++
			if object.LastModified.Before(*oldestObject.LastModified) {
				oldestObject = object
			}
		}
	}

	return snapshotCnt, *oldestObject.Key, nil
}
