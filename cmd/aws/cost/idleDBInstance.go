package cost

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/common"
)

const (
	IdleDBInstanceCheckId                  = "ckia:aws:cost:IdleDBInstance"
	IdleDBInstanceCheckName                = "RDS Idle DB Instances"
	IdleDBInstanceCheckDescription         = "Checks the configuration of your Amazon Relational Database Service (Amazon RDS) for any DB instances that appear to be idle. If a DB instance has not had a connection for a prolonged period of time, you can delete the instance to reduce costs. If persistent storage is needed for data on the instance, you can use lower-cost options such as taking and retaining a DB snapshot. Manually created DB snapshots are retained until you delete them."
	IdleDBInstanceCheckCriteria            = "Any RDS DB instance that has not had a connection in the last 7 days is considered idle."
	IdleDBInstanceCheckRecommendedAction   = "Consider taking a snapshot of the idle DB instance and then either stopping it or deleting it. Stopping the DB instance removes some of the costs for it, but does not remove storage costs. A stopped instance keeps all automated backups based upon the configured retention period. Stopping a DB instance usually incurs additional costs when compared to deleting the instance and then retaining only the final snapshot."
	IdleDBInstanceCheckAdditionalResources = "See comparable AWS Trusted advisor check: https://docs.aws.amazon.com/awssupport/latest/user/cost-optimization-checks.html#amazon-rds-idle-dbs-instances"
)

type IdleDBInstance struct {
	Region                  string `json:"region"`
	DBInstanceName          string `json:"dbInstanceName"`
	MultiAZ                 bool   `json:"multiAZ"`
	InstanceType            string `json:"instanceType"`
	StorageProvisionedInGB  int    `json:"storageProvisionedInGB"`
	DaysSinceLastConnection int    `json:"daysSinceLastConnection"`
	EstimatedMonthlySavings int    `json:"estimatedMonthlySavings"`
}

type IdleDBInstanceCheck struct {
	common.Check
	IdleDBInstances []IdleDBInstance `json:"idleDBInstances"`
}

func FindIdleDBInstances(ctx context.Context, conn client.AWSClient) IdleDBInstanceCheck {
	check := IdleDBInstanceCheck{
		Check: common.Check{
			Id:                  IdleDBInstanceCheckId,
			Name:                IdleDBInstanceCheckName,
			Description:         IdleDBInstanceCheckDescription,
			Criteria:            IdleDBInstanceCheckCriteria,
			RecommendedAction:   IdleDBInstanceCheckRecommendedAction,
			AdditionalResources: IdleDBInstanceCheckAdditionalResources,
		},
	}

	currentTime := time.Now()

	in := &rds.DescribeDBInstancesInput{}
	out, err := conn.RDS.DescribeDBInstances(ctx, in)

	if err != nil {
		fmt.Errorf("Error Listing RDS Instances: %s", err)
	}

	var idleDBInstances []IdleDBInstance
	for _, dbInstance := range out.DBInstances {

		metrics, err := conn.Cloudwatch.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
			MetricName: aws.String("DatabaseConnections"),
			Period:     aws.Int32(3600),
			Namespace:  aws.String("AWS/RDS"),
			Statistics: []types.Statistic{types.StatisticAverage},
			Dimensions: []types.Dimension{
				{
					Name:  aws.String("DBInstanceIdentifier"),
					Value: dbInstance.DBInstanceIdentifier,
				},
			},
			StartTime: aws.Time(currentTime.AddDate(0, 0, -7)),
			EndTime:   aws.Time(currentTime),
		})

		if err != nil {
			return IdleDBInstanceCheck{}
		}

		connectionFound := false
		var idleDBInstance IdleDBInstance
		for _, dataPoint := range metrics.Datapoints {
			if *dataPoint.Average != 0 {
				connectionFound = true
			}
		}

		if !connectionFound {
			// pricingSvc := pricing.NewFromConfig(cfg)
			// filters := []pricingtypes.Filter{
			// 	{
			// 		Field: aws.String("InstanceType"),
			// 		Type:  "TERM_MATCH",
			// 		Value: dbInstance.DBInstanceClass,
			// 	},
			// 	// These two seam to not match what the pricing API is expecting
			// 	{
			// 		Field: aws.String("storage"),
			// 		Type:  "TERM_MATCH",
			// 		Value: dbInstance.StorageType,
			// 	},
			// 	{
			// 		Field: aws.String("databaseEngine"),
			// 		Type:  "TERM_MATCH",
			// 		Value: dbInstance.Engine,
			// 	},
			// 	{
			// 		Field: aws.String("deploymentOption"),
			// 		Type:  "TERM_MATCH",
			// 		Value: aws.String("Single-AZ"),
			// 	},
			// 	{
			// 		Field: aws.String("termType"),
			// 		Type:  "TERM_MATCH",
			// 		Value: aws.String("OnDemand"),
			// 	},
			// 	{
			// 		Field: aws.String("regionCode"),
			// 		Type:  "TERM_MATCH",
			// 		Value: &cfg.Region,
			// 	},
			// 	{
			// 		Field: aws.String("purchaseOption"),
			// 		Type:  "TERM_MATCH",
			// 		Value: aws.String("No Upfront"),
			// 	},
			// }

			// pricingIn := &pricing.GetProductsInput{
			// 	ServiceCode: aws.String("AmazonRDS"),
			// 	Filters:     filters,
			// }
			// pricingData, err := pricingSvc.GetProducts(ctx, pricingIn)

			// if err != nil {
			// 	return IdleDBInstanceCheck{}
			// }

			idleDBInstance.DBInstanceName = aws.ToString(dbInstance.DBInstanceIdentifier)
			idleDBInstance.Region = conn.Region
			idleDBInstance.DaysSinceLastConnection = 7
			idleDBInstance.InstanceType = aws.ToString(dbInstance.DBInstanceClass)
			idleDBInstance.MultiAZ = dbInstance.MultiAZ
			idleDBInstance.StorageProvisionedInGB = int(dbInstance.AllocatedStorage)
			// Still trying to figure out how to get the proper on demand pricing via the API
			// idleDBInstance.EstimatedMonthlySavings = 0
			idleDBInstances = append(idleDBInstances, idleDBInstance)
		}

	}

	check.IdleDBInstances = idleDBInstances
	return check
}
