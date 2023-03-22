package aws

import (
	"github.com/brittandeyoung/ckia/internal/aws/cost"
	"github.com/brittandeyoung/ckia/internal/aws/security"
)

type checkMapping map[string]interface{}

func BuildChecksMap() map[string]interface{} {
	checksMap := checkMapping{
		// Cost Checks go here
		cost.IdleDBInstanceCheckId: cost.FindIdleDBInstances,
		// Security checks go here
		security.RootAccountMFACheckId: security.FindRootAccountsMissingMFA,
	}
	return checksMap
}
