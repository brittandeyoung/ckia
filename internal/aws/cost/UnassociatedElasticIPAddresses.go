package cost

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/common"
)

const (
	UnassociatedElasticIPAddressesCheckId                  = "ckia:aws:cost:UnassociatedElasticIPAddresses"
	UnassociatedElasticIPAddressesCheckName                = "Unassociated Elastic IP Addresses"
	UnassociatedElasticIPAddressesCheckDescription         = "Checks for Elastic IP addresses (EIPs) that are not associated with a running Amazon Elastic Compute Cloud (Amazon EC2) instance. EIPs are static IP addresses designed for dynamic cloud computing. Unlike traditional static IP addresses, EIPs mask the failure of an instance or Availability Zone by remapping a public IP address to another instance in your account. A nominal charge is imposed for an EIP that is not associated with a running instance."
	UnassociatedElasticIPAddressesCheckCriteria            = "An allocated Elastic IP address (EIP) is not associated with a running Amazon EC2 instance."
	UnassociatedElasticIPAddressesCheckRecommendedAction   = "Associate the EIP with a running active instance, or release the unassociated EIP. "
	UnassociatedElasticIPAddressesCheckAdditionalResources = "See comparable AWS Trusted advisor check: https://docs.aws.amazon.com/awssupport/latest/user/cost-optimization-checks.html#unassociated-elastic-ip-addresses"
)

type UnassociatedElasticIPAddresses struct {
	Region    string `json:"region"`
	IpAddress string `json:"ipAddress`
}

type UnassociatedElasticIPAddressesCheck struct {
	common.Check
	UnassociatedElasticIPAddresses []UnassociatedElasticIPAddresses `json:"unassociatedAddresses"`
}

func (v UnassociatedElasticIPAddressesCheck) List() *UnassociatedElasticIPAddressesCheck {
	check := &UnassociatedElasticIPAddressesCheck{
		Check: common.Check{
			Id:                  UnassociatedElasticIPAddressesCheckId,
			Name:                UnassociatedElasticIPAddressesCheckName,
			Description:         UnassociatedElasticIPAddressesCheckDescription,
			Criteria:            UnassociatedElasticIPAddressesCheckCriteria,
			RecommendedAction:   UnassociatedElasticIPAddressesCheckRecommendedAction,
			AdditionalResources: UnassociatedElasticIPAddressesCheckAdditionalResources,
		},
	}
	return check
}

func (v UnassociatedElasticIPAddressesCheck) Run(ctx context.Context, conn client.AWSClient) (*UnassociatedElasticIPAddressesCheck, error) {
	check := new(UnassociatedElasticIPAddressesCheck).List()

	in := &ec2.DescribeAddressesInput{}
	out, err := conn.EC2.DescribeAddresses(ctx, in)

	if err != nil {
		return nil, err
	}

	if len(out.Addresses) == 0 {
		return nil, nil
	}

	var unassociatedAddresses []UnassociatedElasticIPAddresses
	for _, address := range out.Addresses {

		unassociatedAddress := expandUnassociatedAddress(conn, address)

		if unassociatedAddress.IpAddress != "" {
			unassociatedAddresses = append(unassociatedAddresses, unassociatedAddress)

		}
	}

	check.UnassociatedElasticIPAddresses = unassociatedAddresses
	return check, nil
}

func expandUnassociatedAddress(conn client.AWSClient, address types.Address) UnassociatedElasticIPAddresses {
	var unassociatedAddress UnassociatedElasticIPAddresses
	if address.AssociationId == nil {
		unassociatedAddress.Region = conn.Region
		unassociatedAddress.IpAddress = aws.ToString(address.PublicIp)
	}
	return unassociatedAddress
}
