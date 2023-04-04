package cost

import (
	"context"
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

type UnderutilizedEBSVolume struct {
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
	UnderutilizedEBSVolumes []UnderutilizedEBSVolume `json:"underutilizedVolumes"`
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

func (v UnderutilizedEBSVolumesCheck) Run(ctx context.Context, conn client.AWSClient) (*UnderutilizedEBSVolumesCheck, error) {
	check := new(UnderutilizedEBSVolumesCheck).List()

	currentTime := time.Now()

	in := &ec2.DescribeVolumesInput{}
	var volumes []types.Volume

	paginator := ec2.NewDescribeVolumesPaginator(conn.EC2, in, func(o *ec2.DescribeVolumesPaginatorOptions) {})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)

		if err != nil {
			return nil, err
		}
		volumes = append(volumes, output.Volumes...)

	}

	if len(volumes) == 0 {
		return nil, nil
	}

	var underutilizedVolumes []UnderutilizedEBSVolume
	for _, volume := range volumes {

		var underutilizedVolume UnderutilizedEBSVolume

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
			return nil, err
		}

		underutilizedVolume = expandUnderutilizedVolume(conn, volume, metrics.Datapoints)

		if underutilizedVolume.SnapshotId != "" {
			snapshots, err := conn.EC2.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
				SnapshotIds: []string{underutilizedVolume.SnapshotId},
			})

			if err != nil {
				return nil, err
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
	return check, nil
}

func expandUnderutilizedVolume(conn client.AWSClient, volume types.Volume, dataPoints []cloudWatchTypes.Datapoint) UnderutilizedEBSVolume {
	var underutilizedVolume UnderutilizedEBSVolume
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

func expandSnapshot(snapshots []types.Snapshot, volume UnderutilizedEBSVolume) UnderutilizedEBSVolume {
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
