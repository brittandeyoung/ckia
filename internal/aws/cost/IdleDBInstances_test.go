package cost

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func TestExpandConnections_basic(t *testing.T) {
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

	daysSinceConnection, connectionFound := expandConnections(dataPoints)

	if connectionFound || daysSinceConnection != 14 {
		t.Fatal(`Connection found when no connection present`)
	}

}

func TestExpandConnections_withConnectionsInSevenDays(t *testing.T) {
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
			Timestamp: aws.Time(time.Now().Add(time.Hour * -48)),
			Average:   aws.Float64(1.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now().Add(time.Hour * -72)),
			Average:   aws.Float64(1.0),
			Unit:      "Count",
		},
	}

	daysSinceConnection, connectionFound := expandConnections(dataPoints)

	if !connectionFound {
		t.Fatal(`Connection not found when connection is present`)
	}

	if daysSinceConnection != 2 {
		t.Fatalf(`Days Since Connetion should be 3, Got %d`, daysSinceConnection)
	}
}

func TestExpandConnections_withConnectionsAfterSevenDays(t *testing.T) {
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
			Timestamp: aws.Time(time.Now().Add(time.Hour * -192)),
			Average:   aws.Float64(1.0),
			Unit:      "Count",
		},
		{
			Timestamp: aws.Time(time.Now().Add(time.Hour * -240)),
			Average:   aws.Float64(1.0),
			Unit:      "Count",
		},
	}

	daysSinceConnection, connectionFound := expandConnections(dataPoints)

	if connectionFound {
		t.Fatal(`Connection reported within the last 7 days when not present.`)
	}

	if daysSinceConnection != 8 {
		t.Fatalf(`Days Since Connetion should be 8, Got %d`, daysSinceConnection)
	}
}
