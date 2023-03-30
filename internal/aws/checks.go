package aws

import (
	"github.com/brittandeyoung/ckia/internal/aws/cost"
	"github.com/brittandeyoung/ckia/internal/aws/security"
)

type checkMapping map[string]interface{}

func BuildChecksMap() map[string]interface{} {
	checksMap := checkMapping{
		// Cost Checks go here
		cost.IdleDBInstanceCheckId:                 new(cost.IdleDBInstanceCheck),
		cost.UnderutilizedEBSVolumesCheckId:        new(cost.UnderutilizedEBSVolumesCheck),
		cost.UnassociatedElasticIPAddressesCheckId: new(cost.UnassociatedElasticIPAddressesCheck),
		// Security checks go here
		security.RootAccountMissingMFACheckId: new(security.RootAccountMissingMFACheck),
	}
	return checksMap
}
