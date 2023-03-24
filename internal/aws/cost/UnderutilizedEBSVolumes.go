package cost

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudWatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/common"
)

const (
	UnderutilizedEBSVolumesCheckId                  = "ckia:aws:cost:UnderutilizedEBSVolumes"
	UnderutilizedEBSVolumesCheckName                = "Underutilized Amazon EBS Volumes"
	UnderutilizedEBSVolumesCheckDescription         = "Checks Amazon Elastic Block Store (Amazon EBS) volume configurations and warns when volumes appear to be underutilized. Charges begin when a volume is created. If a volume remains unattached or has very low write activity (excluding boot volumes) for a period of time, the volume is underutilized. We recommend that you remove underutilized volumes to reduce costs."
	UnderutilizedEBSVolumesCheckCriteria            = "A volume is unattached or had less than 1 IOPS per day for the past 7 days."
	UnderutilizedEBSVolumesCheckRecommendedAction   = "Consider creating a snapshot and deleting the volume to reduce costs."
	UnderutilizedEBSVolumesCheckAdditionalResources = "See comparable AWS Trusted advisor check: https://docs.aws.amazon.com/awssupport/latest/user/cost-optimization-checks.html#underutilized-amazon-ebs-volumes"
)

type UnderutilizedEBSVolumes struct {
	Region             string `json:"region"`
	VolumeId           string `json:"volumeId"`
	VolumeName         string `json:"volumeName"`
	VolumeType         string `json:"volumeType"`
	VolumeSize         int    `json:"volumeSize"`
	MonthlyStorageCost int    `json:"monthlyStorageCost"`
	SnapshotId         string `json:"snapshotId"`
	SnapshotName       string `json:"snapshotName"`
	SnapshotAge        int    `json:"snapshotAge"`
}

type UnderutilizedEBSVolumesCheck struct {
	common.Check
	UnderutilizedEBSVolumes []UnderutilizedEBSVolumes `json:"underutilizedVolumes"`
}

func (v UnderutilizedEBSVolumesCheck) List() *UnderutilizedEBSVolumesCheck {
	check := &UnderutilizedEBSVolumesCheck{
		Check: common.Check{
			Id:                  UnderutilizedEBSVolumesCheckId,
			Name:                UnderutilizedEBSVolumesCheckName,
			Description:         UnderutilizedEBSVolumesCheckDescription,
			Criteria:            UnderutilizedEBSVolumesCheckCriteria,
			RecommendedAction:   UnderutilizedEBSVolumesCheckRecommendedAction,
			AdditionalResources: UnderutilizedEBSVolumesCheckAdditionalResources,
		},
	}
	return check
}

func (v UnderutilizedEBSVolumesCheck) Run(ctx context.Context, conn client.AWSClient) *UnderutilizedEBSVolumesCheck {
	check := new(UnderutilizedEBSVolumesCheck).List()

	currentTime := time.Now()

	in := &ec2.DescribeVolumesInput{}
	out, err := conn.EC2.DescribeVolumes(ctx, in)

	if err != nil {
		fmt.Errorf("Error Listing Volumes: %s", err)
	}

	if len(out.Volumes) == 0 {
		return nil
	}

	var underutilizedVolumes []UnderutilizedEBSVolumes
	for _, volume := range out.Volumes {

		var underutilizedVolume UnderutilizedEBSVolumes

		metrics, err := conn.Cloudwatch.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
			MetricName: aws.String("VolumeReadOps"),
			Period:     aws.Int32(3600),
			Namespace:  aws.String("AWS/EBS"),
			Statistics: []cloudWatchTypes.Statistic{cloudWatchTypes.StatisticSum},
			Dimensions: []cloudWatchTypes.Dimension{
				{
					Name:  aws.String("DBInstanceIdentifier"),
					Value: volume.VolumeId,
				},
			},
			StartTime: aws.Time(currentTime.AddDate(0, 0, -14)),
			EndTime:   aws.Time(currentTime),
		})

		if err != nil {
			fmt.Errorf("Error Retrieving EBS Metrics for Volume: %s", aws.ToString(volume.VolumeId))
		}

		underutilizedVolume = expandUnderutilizedVolume(conn, volume, metrics.Datapoints)

		if underutilizedVolume.SnapshotId != "" {
			snapshots, err := conn.EC2.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
				SnapshotIds: []string{underutilizedVolume.SnapshotId},
			})

			if err != nil {
				fmt.Errorf("Error Listing Snapshots: %s", err)
			}

			underutilizedVolume = expandSnapshot(snapshots.Snapshots, underutilizedVolume)
			// Still trying to figure out how to get the proper pricing via the API
			// underutilizedVolume.MonthlyStorageCost = 0
		}
		if underutilizedVolume.VolumeId != "" {
			underutilizedVolumes = append(underutilizedVolumes, underutilizedVolume)

		}
	}

	check.UnderutilizedEBSVolumes = underutilizedVolumes
	return check
}

func expandUnderutilizedVolume(conn client.AWSClient, volume types.Volume, dataPoints []cloudWatchTypes.Datapoint) UnderutilizedEBSVolumes {
	var underutilizedVolume UnderutilizedEBSVolumes
	iopsFound := false
	for _, dataPoint := range dataPoints {
		if aws.ToFloat64(dataPoint.Average) != 0 {
			iopsFound = true
		}
	}
	if volume.State == types.VolumeStateAvailable && !iopsFound {
		underutilizedVolume.Region = conn.Region
		underutilizedVolume.VolumeId = aws.ToString(volume.VolumeId)
		for _, tag := range volume.Tags {
			if aws.ToString(tag.Key) == "Name" {
				underutilizedVolume.VolumeName = aws.ToString(tag.Value)
			}
		}
		underutilizedVolume.VolumeType = aws.ToString((*string)(&volume.VolumeType))
		underutilizedVolume.VolumeSize = int(aws.ToInt32(volume.Size))
		underutilizedVolume.SnapshotId = aws.ToString(volume.SnapshotId)
	}
	return underutilizedVolume
}

func expandSnapshot(snapshots []types.Snapshot, volume UnderutilizedEBSVolumes) UnderutilizedEBSVolumes {
	if len(snapshots) > 0 {
		snapshot := snapshots[0]
		duration := time.Since(aws.ToTime(snapshot.StartTime))
		volume.SnapshotAge = int(duration.Hours() / 24)
		for _, tag := range snapshot.Tags {
			if aws.ToString(tag.Key) == "Name" {
				volume.SnapshotName = aws.ToString(tag.Value)
			}
		}
	}
	return volume
}
