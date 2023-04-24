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
	RootAccountMissingMFACheckId                  = "ckia:aws:security:RootAccountMissingMFA"
	RootAccountMissingMFACheckName                = "MFA on Root Account"
	RootAccountMissingMFACheckDescription         = "Checks the root account and warns if multi-factor authentication (MFA) is not enabled. For increased security, we recommend that you protect your account by using MFA, which requires a user to enter a unique authentication code from their MFA hardware or virtual device when interacting with the AWS Management Console and associated websites."
	RootAccountMissingMFACheckCriteria            = "MFA is not enabled on the root account."
	RootAccountMissingMFACheckRecommendedAction   = "Log in to your root account and activate an MFA device. "
	RootAccountMissingMFACheckAdditionalResources = "Using Multi-Factor Authentication (MFA) Devices with AWS: https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_ManagingMFA.html"
)

type RootAccountMissingMFA struct {
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`
}

type RootAccountMissingMFACheck struct {
	common.Check
	AccountId              string                  `json:"accountId"`
	RootAccountsMissingMFA []RootAccountMissingMFA `json:"rootAccountsMissingMFA"`
}

func (v *RootAccountMissingMFACheck) List() *RootAccountMissingMFACheck {
	v.Check = common.Check{
		Id:                  RootAccountMissingMFACheckId,
		Name:                RootAccountMissingMFACheckName,
		Description:         RootAccountMissingMFACheckDescription,
		Criteria:            RootAccountMissingMFACheckCriteria,
		RecommendedAction:   RootAccountMissingMFACheckRecommendedAction,
		AdditionalResources: RootAccountMissingMFACheckAdditionalResources,
	}

	return v
}

func (v *RootAccountMissingMFACheck) Run(ctx context.Context, conn client.AWSClient) (*RootAccountMissingMFACheck, error) {
	v = v.List()

	accountSummary, err := conn.IAM.GetAccountSummary(ctx, &iam.GetAccountSummaryInput{})

	if err != nil {
		return nil, err
	}

	identity, err := conn.STS.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, err
	}

	account := expandRootAccountMissingMFA(accountSummary.SummaryMap, aws.ToString(identity.Account))

	if account != (RootAccountMissingMFA{}) {
		v.RootAccountsMissingMFA = []RootAccountMissingMFA{account}
	}

	return v, nil
}

func expandRootAccountMissingMFA(summaryMap map[string]int32, accountNumber string) RootAccountMissingMFA {
	var rootAccountMissingMFA RootAccountMissingMFA
	if summaryMap["AccountMFAEnabled"] != 1 {
		rootAccountMissingMFA.AccountId = accountNumber
		// Can only get account name from the organization account
		// check.AccountName
	}
	return rootAccountMissingMFA
}
