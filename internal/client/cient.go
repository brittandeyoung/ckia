package client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSClient struct {
	Cloudwatch *cloudwatch.Client
	EC2        *ec2.Client
	ELBv2      *elasticloadbalancingv2.Client
	IAM        *iam.Client
	Pricing    *pricing.Client
	RDS        *rds.Client
	Region     string
	STS        *sts.Client
}

func InitiateClient(cfg aws.Config) AWSClient {
	client := AWSClient{
		Cloudwatch: cloudwatch.NewFromConfig(cfg),
		EC2:        ec2.NewFromConfig(cfg),
		ELBv2:      elasticloadbalancingv2.NewFromConfig(cfg),
		IAM:        iam.NewFromConfig(cfg),
		Pricing:    pricing.NewFromConfig(cfg),
		RDS:        rds.NewFromConfig(cfg),
		Region:     cfg.Region,
		STS:        sts.NewFromConfig(cfg),
	}

	return client
}
