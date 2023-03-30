package cost

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/create"
)

func TestExpandUnderutilizedVolume_basic(t *testing.T) {
	dataPoints := []types.Datapoint{
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
	}
	volume := ec2Types.Volume{
		Attachments: []ec2Types.VolumeAttachment{},
		Size:        aws.Int32(20),
		SnapshotId:  aws.String("snap-0240fe3027dd6b4wa0"),
		State:       ec2Types.VolumeStateAvailable,
		VolumeId:    aws.String("vol-02e71c945942481e85"),
		VolumeType:  ec2Types.VolumeTypeGp2,
		Tags: []ec2Types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("MyVolumeName"),
			},
		},
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	conn := client.InitiateClient(cfg)
	underutilizedVolume := expandUnderutilizedVolume(conn, volume, dataPoints)

	if underutilizedVolume == (UnderutilizedEBSVolumes{}) {
		create.TestFailureNonEmptyStruct(t)
	}
	if underutilizedVolume.Region != "us-east-1" {
		t.Fatal(`Volume Region did not set properly from expand function.`)
	}
	if underutilizedVolume.VolumeId != "vol-02e71c945942481e85" {
		t.Fatal(`Volume Region did not set properly from expand function.`)
	}
	if underutilizedVolume.VolumeName != "MyVolumeName" {
		t.Fatal(`Volume Region did not set properly from expand function.`)
	}
	if underutilizedVolume.VolumeType != "gp2" {
		t.Fatal(`Volume Region did not set properly from expand function.`)
	}
	if underutilizedVolume.VolumeSize != 20 {
		t.Fatal(`Volume Region did not set properly from expand function.`)
	}
	if underutilizedVolume.SnapshotId != "snap-0240fe3027dd6b4wa0" {
		t.Fatal(`Volume Region did not set properly from expand function.`)
	}
}

func TestExpandUnderutilizedVolume_attachedVolume(t *testing.T) {
	dataPoints := []types.Datapoint{
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
	}
	volume := ec2Types.Volume{
		Attachments: []ec2Types.VolumeAttachment{},
		Size:        aws.Int32(20),
		SnapshotId:  aws.String("snap-0240fe3027dd6b4wa0"),
		State:       ec2Types.VolumeStateInUse,
		VolumeId:    aws.String("vol-02e71c945942481e85"),
		VolumeType:  ec2Types.VolumeTypeGp2,
		Tags: []ec2Types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("MyVolumeName"),
			},
		},
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	conn := client.InitiateClient(cfg)
	underutilizedVolume := expandUnderutilizedVolume(conn, volume, dataPoints)

	if underutilizedVolume != (UnderutilizedEBSVolumes{}) {
		create.TestFailureNonEmptyStruct(t)
	}
}

func TestExpandUnderutilizedVolume_activeVolume(t *testing.T) {
	dataPoints := []types.Datapoint{
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(5.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(1.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(6.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now()),
			Average:   aws.Float64(0.0),
			Unit:      "Count",
		},
	}
	volume := ec2Types.Volume{
		Attachments: []ec2Types.VolumeAttachment{},
		Size:        aws.Int32(20),
		SnapshotId:  aws.String("snap-0240fe3027dd6b4wa0"),
		State:       ec2Types.VolumeStateAvailable,
		VolumeId:    aws.String("vol-02e71c945942481e85"),
		VolumeType:  ec2Types.VolumeTypeGp2,
		Tags: []ec2Types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("MyVolumeName"),
			},
		},
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	conn := client.InitiateClient(cfg)
	underutilizedVolume := expandUnderutilizedVolume(conn, volume, dataPoints)

	if underutilizedVolume != (UnderutilizedEBSVolumes{}) {
		create.TestFailureNonEmptyStruct(t)
	}
}

func TestExpandSnapshot_basic(t *testing.T) {
	snapshots := []ec2Types.Snapshot{
		{
			SnapshotId: aws.String("snap-0240fe3027dd6b4wa0"),
			StartTime:  aws.Time(time.Now().Add(time.Hour * -192)),
			Tags: []ec2Types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("MySnapshotName"),
				},
			},
		},
	}
	var underutilizedVolume UnderutilizedEBSVolumes
	underutilizedVolume = expandSnapshot(snapshots, underutilizedVolume)

	if underutilizedVolume.SnapshotAge != 8 {
		t.Fatal(`Snapshot Age is not computed or set correctly by expand function.`)
	}
	if underutilizedVolume.SnapshotName != "MySnapshotName" {
		t.Fatal(`Snapshot Name is not set correctly by expand function.`)
	}

}
