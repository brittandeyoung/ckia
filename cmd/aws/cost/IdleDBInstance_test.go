package cost

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func TestExpandConnections(t *testing.T) {
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

	connections, connectionFound := expandConnections(dataPoints)

	if connectionFound || len(connections) > 0 {
		t.Fatal(`Connection found when no connection present`)
	}

}
