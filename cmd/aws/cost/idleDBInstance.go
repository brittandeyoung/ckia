package cost

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/brittandeyoung/ckia/internal/common"
)

const (
	IdleDBInstanceCheckId                  = "ckia:aws:cost:IdleDBInstanceCheck"
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

func ListRDSIdleDB() IdleDBInstanceCheck {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

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
	svc := rds.NewFromConfig(cfg)
	cloudwatchSvc := cloudwatch.NewFromConfig(cfg)

	currentTime := time.Now()

	in := &rds.DescribeDBInstancesInput{}

	out, err := svc.DescribeDBInstances(context.TODO(), in)

	if err != nil {
		fmt.Errorf("Error Listing RDS Instances: %s", err)
	}

	var idleDBInstances []IdleDBInstance
	for _, dbInstance := range out.DBInstances {

		metrics, err := cloudwatchSvc.GetMetricStatistics(context.TODO(), &cloudwatch.GetMetricStatisticsInput{
			MetricName: aws.String("DatabaseConnections"),
			Period:     aws.Int32(3600),
			Namespace:  aws.String("AWS/RDS"),
			Statistics: []types.Statistic{types.StatisticAverage},
			Dimensions: []types.Dimension{
				types.Dimension{
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
		for _, datapoint := range metrics.Datapoints {
			if *datapoint.Average == 0 {
				continue
			} else {
				connectionFound = true
			}
		}

		if connectionFound {
			// pricingSvc := pricing.NewFromConfig(cfg)
			// filters := []pricingtypes.Filter{
			// 	pricingtypes.Filter{
			// 		Field: aws.String("InstanceType"),
			// 		Type:  "TERM_MATCH",
			// 		Value: dbInstance.DBInstanceClass,
			// 	},
			// 	pricingtypes.Filter{
			// 		Field: aws.String("InstanceType"),
			// 		Type:  "TERM_MATCH",
			// 		Value: dbInstance.DBInstanceClass,
			// 	},
			// }

			// in = &pricing.GetProductsInput{
			// 	ServiceCode: "AmazonRDS",
			// 	Filters:     filters,
			// }
			// pricingSvc.GetProducts()
			idleDBInstance.DBInstanceName = aws.ToString(dbInstance.DBInstanceIdentifier)
			idleDBInstance.Region = cfg.Region
			idleDBInstance.DaysSinceLastConnection = 7
			idleDBInstance.InstanceType = aws.ToString(dbInstance.DBInstanceClass)
			idleDBInstance.MultiAZ = dbInstance.MultiAZ
			idleDBInstance.StorageProvisionedInGB = int(dbInstance.AllocatedStorage)
			idleDBInstance.EstimatedMonthlySavings = 100
			idleDBInstances = append(idleDBInstances, idleDBInstance)
		}

	}

	check.IdleDBInstances = idleDBInstances
	return check
}
