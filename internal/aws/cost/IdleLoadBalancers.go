package cost

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	lbTypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/common"
)

const (
	IdleLoadBalancersCheckId                  = "ckia:aws:cost:IdleLoadBalancers"
	IdleLoadBalancersCheckName                = "Idle Load Balancers"
	IdleLoadBalancersCheckDescription         = "Checks your Elastic Load Balancing configuration for load balancers that are idle. Any load balancer that is configured accrues charges. If a load balancer has no associated back-end instances, or if network traffic is severely limited, the load balancer is not being used effectively. This check currently only checks for Classic Load Balancer type within ELB service. It does not include other ELB types (Application Load Balancer, Network Load Balancer)"
	IdleLoadBalancersCheckCriteria            = "A load balancer has no active back-end instances. A load balancer has no healthy back-end instances. A load balancer has had less than 100 requests per day for the last 7 days."
	IdleLoadBalancersCheckRecommendedAction   = "If your load balancer has no active back-end instances, consider registering instances or deleting your load balancer. If your load balancer has no healthy back-end instances, troubleshoot why they are un healthy or evaluate for removal. If your load balancer has had a low request count, consider deleting your load balancer. See Delete Your Load Balancer."
	IdleLoadBalancersCheckAdditionalResources = "See comparable AWS Trusted advisor check: https://docs.aws.amazon.com/awssupport/latest/user/cost-optimization-checks.html#idle-load-balancers"

	IdleLoadBalancerReasonNoActiveInstances  = "no active back-end instances"
	IdleLoadBalancerReasonNoHealthyInstances = "no healthy back-end instances"
	IdleLoadBalancerReasonLowRequestCount    = "low request count"
)

type IdleLoadBalancer struct {
	Region                  string `json:"region"`
	LoadBalancerName        string `json:"loadBalancerName"`
	Reason                  string `json:"reason"`
	EstimatedMonthlySavings int    `json:"estimatedMonthlySavings"`
}

type IdleLoadBalancersCheck struct {
	common.Check
	IdleLoadBalancers []IdleLoadBalancer `json:"idleLoadBalancers"`
}

func (v IdleLoadBalancersCheck) List() *IdleLoadBalancersCheck {
	check := &IdleLoadBalancersCheck{
		Check: common.Check{
			Id:                  IdleLoadBalancersCheckId,
			Name:                IdleLoadBalancersCheckName,
			Description:         IdleLoadBalancersCheckDescription,
			Criteria:            IdleLoadBalancersCheckCriteria,
			RecommendedAction:   IdleLoadBalancersCheckRecommendedAction,
			AdditionalResources: IdleLoadBalancersCheckAdditionalResources,
		},
	}
	return check
}

func (v IdleLoadBalancersCheck) Run(ctx context.Context, conn client.AWSClient) (*IdleLoadBalancersCheck, error) {
	check := new(IdleLoadBalancersCheck).List()

	currentTime := time.Now()
	var loadBalancers []lbTypes.LoadBalancer
	in := &elasticloadbalancingv2.DescribeLoadBalancersInput{}

	paginator := elasticloadbalancingv2.NewDescribeLoadBalancersPaginator(conn.ELBv2, in, func(o *elasticloadbalancingv2.DescribeLoadBalancersPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)

		if err != nil {
			return nil, err
		}
		loadBalancers = append(loadBalancers, output.LoadBalancers...)

	}

	if len(loadBalancers) == 0 {
		return nil, nil
	}

	var idleLoadBalancers []IdleLoadBalancer
	for _, lb := range loadBalancers {

		targetGroups, err := conn.ELBv2.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{
			LoadBalancerArn: lb.LoadBalancerArn,
		})

		if err != nil {
			return nil, err
		}

		lbIsIdle := false
		var idleLoadBalancer IdleLoadBalancer

		for _, group := range targetGroups.TargetGroups {
			if !lbIsIdle {
				health, err := conn.ELBv2.DescribeTargetHealth(ctx, &elasticloadbalancingv2.DescribeTargetHealthInput{
					TargetGroupArn: group.TargetGroupArn,
				})

				if err != nil {
					return nil, err
				}

				idleLoadBalancer, lbIsIdle = expandInactiveLoadBalancer(idleLoadBalancer, health.TargetHealthDescriptions)

				if !lbIsIdle {
					idleLoadBalancer, lbIsIdle = expandUnhealthyLoadBalancer(idleLoadBalancer, health.TargetHealthDescriptions)
				}
			}
		}

		if !lbIsIdle {
			metrics, err := conn.Cloudwatch.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
				MetricName: aws.String("RequestCount"),
				Period:     aws.Int32(3600),
				Namespace:  aws.String("AWS/ApplicationELB"),
				Statistics: []types.Statistic{types.StatisticAverage},
				Dimensions: []types.Dimension{
					{
						Name:  aws.String("LoadBalancer"),
						Value: lb.LoadBalancerName,
					},
				},
				StartTime: aws.Time(currentTime.AddDate(0, 0, -7)),
				EndTime:   aws.Time(currentTime),
			})

			if err != nil {
				return nil, err
			}

			idleLoadBalancer, lbIsIdle = expandLowRequestCountLoadBalancer(idleLoadBalancer, metrics.Datapoints)
		}

		if lbIsIdle {
			idleLoadBalancer.LoadBalancerName = aws.ToString(lb.LoadBalancerName)
			idleLoadBalancer.Region = conn.Region
			// Still trying to figure out how to get the proper on demand pricing via the API
			// idleLoadBalancer.EstimatedMonthlySavings = 0
			idleLoadBalancers = append(idleLoadBalancers, idleLoadBalancer)
		}

	}

	check.IdleLoadBalancers = idleLoadBalancers
	return check, nil
}

func expandInactiveLoadBalancer(idleLoadBalancer IdleLoadBalancer, descriptions []lbTypes.TargetHealthDescription) (IdleLoadBalancer, bool) {
	if len(descriptions) == 0 {
		idleLoadBalancer.Reason = IdleLoadBalancerReasonNoActiveInstances
		return idleLoadBalancer, true
	}

	return idleLoadBalancer, false
}

func expandUnhealthyLoadBalancer(idleLoadBalancer IdleLoadBalancer, descriptions []lbTypes.TargetHealthDescription) (IdleLoadBalancer, bool) {
	if len(descriptions) > 0 {
		for _, description := range descriptions {
			if description.TargetHealth.State != lbTypes.TargetHealthStateEnumUnhealthy {
				return idleLoadBalancer, false
			}
		}
	}
	idleLoadBalancer.Reason = IdleLoadBalancerReasonNoHealthyInstances
	return idleLoadBalancer, true
}

func expandLowRequestCountLoadBalancer(idleLoadBalancer IdleLoadBalancer, dataPoints []types.Datapoint) (IdleLoadBalancer, bool) {
	for _, dataPoint := range dataPoints {
		if aws.ToFloat64(dataPoint.Sum) > 100 {
			return idleLoadBalancer, false
		}
	}

	idleLoadBalancer.Reason = IdleLoadBalancerReasonLowRequestCount
	return idleLoadBalancer, true
}
