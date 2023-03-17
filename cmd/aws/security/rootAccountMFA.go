package security

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/common"
)

const (
	rootAccountMFACheckId                  = "ckia:aws:cost:rootAccountMFACheck"
	rootAccountMFACheckName                = "MFA on Root Account"
	rootAccountMFACheckDescription         = "Checks the root account and warns if multi-factor authentication (MFA) is not enabled. For increased security, we recommend that you protect your account by using MFA, which requires a user to enter a unique authentication code from their MFA hardware or virtual device when interacting with the AWS Management Console and associated websites."
	rootAccountMFACheckCriteria            = "MFA is not enabled on the root account."
	rootAccountMFACheckRecommendedAction   = "Log in to your root account and activate an MFA device. "
	rootAccountMFACheckAdditionalResources = "Using Multi-Factor Authentication (MFA) Devices with AWS: https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_ManagingMFA.html"
)

type RootAccountMFA struct {
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`
}

type RootAccountMFACheck struct {
	common.Check
	AccountId              string           `json:"accountId"`
	RootAccountsMissingMFA []RootAccountMFA `json:"rootAccountsMissingMFA"`
}

func FindRootAccountsMissingMFA(ctx context.Context, conn client.AWSClient) RootAccountMFACheck {
	check := RootAccountMFACheck{
		Check: common.Check{
			Id:                  rootAccountMFACheckId,
			Name:                rootAccountMFACheckName,
			Description:         rootAccountMFACheckDescription,
			Criteria:            rootAccountMFACheckCriteria,
			RecommendedAction:   rootAccountMFACheckRecommendedAction,
			AdditionalResources: rootAccountMFACheckAdditionalResources,
		},
	}

	accountSummary, err := conn.IAM.GetAccountSummary(ctx, &iam.GetAccountSummaryInput{})

	if err != nil {
		return RootAccountMFACheck{}
	}

	if accountSummary.SummaryMap["AccountMFAEnabled"] != 1 {
		identity, err := conn.STS.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})

		if err != nil {
			return RootAccountMFACheck{}
		}

		accounts := []RootAccountMFA{
			{
				AccountId: aws.ToString(identity.Account),
				// Can only get account name from the organization account
				// check.AccountName
			},
		}

		check.RootAccountsMissingMFA = accounts

		// Can only get account name from the organization account
		// check.AccountName =
	}

	return check
}
