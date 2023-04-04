package cost

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/brittandeyoung/ckia/internal/create"
)

func TestExpandInactiveLoadBalancer_basic(t *testing.T) {
	idleLoadBalancer := IdleLoadBalancer{}
	descriptions := []types.TargetHealthDescription{}
	lb, lbIsIdle := expandInactiveLoadBalancer(idleLoadBalancer, descriptions)

	if lb == (IdleLoadBalancer{}) {
		create.TestFailureEmptyStruct(t)
	}

	if lb.Reason != IdleLoadBalancerReasonNoActiveInstances {
		create.TestFailureAttribute(t, "Reason", IdleLoadBalancerReasonNoActiveInstances)
	}

	if !lbIsIdle {
		create.TestFailureAttribute(t, "lbIsIdle", "true")
	}
}

func TestExpandUnhealthyLoadBalancer_basic(t *testing.T) {
	idleLoadBalancer := IdleLoadBalancer{}
	descriptions := []types.TargetHealthDescription{
		{
			TargetHealth: &types.TargetHealth{
				State: types.TargetHealthStateEnumUnhealthy,
			},
		},
		{
			TargetHealth: &types.TargetHealth{
				State: types.TargetHealthStateEnumUnhealthy,
			},
		},
		{
			TargetHealth: &types.TargetHealth{
				State: types.TargetHealthStateEnumUnhealthy,
			},
		},
	}
	lb, lbIsIdle := expandUnhealthyLoadBalancer(idleLoadBalancer, descriptions)

	if lb == (IdleLoadBalancer{}) {
		create.TestFailureEmptyStruct(t)
	}

	if lb.Reason != IdleLoadBalancerReasonNoHealthyInstances {
		create.TestFailureAttribute(t, "Reason", IdleLoadBalancerReasonNoHealthyInstances)
	}

	if !lbIsIdle {
		create.TestFailureAttribute(t, "lbIsIdle", "true")
	}
}

func TestExpandLowRequestCountLoadBalancer_basic(t *testing.T) {
	idleLoadBalancer := IdleLoadBalancer{}
	dataPoints := []cloudwatchTypes.Datapoint{
		{
			Sum: aws.Float64(80.0),
		},
		{
			Sum: aws.Float64(80.0),
		},
		{
			Sum: aws.Float64(80.0),
		},
	}

	lb, lbIsIdle := expandLowRequestCountLoadBalancer(idleLoadBalancer, dataPoints)

	if lb == (IdleLoadBalancer{}) {
		create.TestFailureEmptyStruct(t)
	}

	if lb.Reason != IdleLoadBalancerReasonLowRequestCount {
		create.TestFailureAttribute(t, "Reason", IdleLoadBalancerReasonLowRequestCount)
	}

	if !lbIsIdle {
		create.TestFailureAttribute(t, "lbIsIdle", "true")
	}
}

func TestExpandLowRequestCountLoadBalancer_notIdle(t *testing.T) {
	idleLoadBalancer := IdleLoadBalancer{}
	dataPoints := []cloudwatchTypes.Datapoint{
		{
			Sum: aws.Float64(120.0),
		},
		{
			Sum: aws.Float64(80.0),
		},
		{
			Sum: aws.Float64(80.0),
		},
	}

	lb, lbIsIdle := expandLowRequestCountLoadBalancer(idleLoadBalancer, dataPoints)

	if lb != (IdleLoadBalancer{}) {
		create.TestFailureNonEmptyStruct(t)
	}

	if lbIsIdle {
		create.TestFailureAttribute(t, "lbIsIdle", "false")
	}
}
