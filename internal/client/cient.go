package client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSClient struct {
	Cloudwatch *cloudwatch.Client
	IAM        *iam.Client
	Pricing    *pricing.Client
	RDS        *rds.Client
	Region     string
	STS        *sts.Client
}

func InitiateClient(cfg aws.Config) AWSClient {
	client := AWSClient{
		Cloudwatch: cloudwatch.NewFromConfig(cfg),
		IAM:        iam.NewFromConfig(cfg),
		Pricing:    pricing.NewFromConfig(cfg),
		RDS:        rds.NewFromConfig(cfg),
		Region:     cfg.Region,
	}

	return client
}
