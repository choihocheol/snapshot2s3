package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client *s3.Client
	bucket string
}
