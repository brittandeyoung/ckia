package cost

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/create"
)

func TestExpandUnassociatedAddress_basic(t *testing.T) {
	address := types.Address{
		PublicIp:           aws.String("18.214.64.132"),
		AllocationId:       nil,
		Domain:             types.DomainType("vpc"),
		PublicIpv4Pool:     aws.String("amazon"),
		NetworkBorderGroup: aws.String("us-east-1"),
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	conn := client.InitiateClient(cfg)

	unassociatedAddress := expandUnassociatedAddress(conn, address)

	if unassociatedAddress == (UnassociatedElasticIPAddress{}) {
		create.TestFailureEmptyStruct(t)
	}

	if unassociatedAddress.IPAddress != "18.214.64.132" {
		t.Fatalf(`IP Address field not set properly.`)
	}

	if unassociatedAddress.Region != "us-east-1" {
		t.Fatalf(`Region field not set properly.`)
	}
}

func TestExpandUnassociatedAddress_none(t *testing.T) {
	address := types.Address{
		PublicIp:           aws.String("18.214.64.132"),
		AllocationId:       aws.String("eipalloc-0287c07cca688eb9a"),
		AssociationId:      aws.String("eipassoc-01923845827937a"),
		Domain:             types.DomainType("vpc"),
		PublicIpv4Pool:     aws.String("amazon"),
		NetworkBorderGroup: aws.String("us-east-1"),
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	conn := client.InitiateClient(cfg)

	unassociatedAddress := expandUnassociatedAddress(conn, address)

	if unassociatedAddress != (UnassociatedElasticIPAddress{}) {
		create.TestFailureNonEmptyStruct(t)
	}
}
