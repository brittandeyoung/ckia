package cost

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/common"
)

const (
	HighErrorRateLambdaFunctionsCheckId                  = "ckia:aws:cost:HighErrorRateLambdaFunctions"
	HighErrorRateLambdaFunctionsCheckName                = "Lambda Functions with Excessive Timeouts"
	HighErrorRateLambdaFunctionsCheckDescription         = "Checks for Lambda functions with high error rates that might result in higher costs."
	HighErrorRateLambdaFunctionsCheckCriteria            = "Functions where > 10% of invocations end in error on any given day within the last 7 days."
	HighErrorRateLambdaFunctionsCheckRecommendedAction   = "Consider the following guidelines to reduce errors. Function errors include errors returned by the function's code and errors returned by the function's runtime."
	HighErrorRateLambdaFunctionsCheckAdditionalResources = "See comparable AWS Trusted advisor check: https://docs.aws.amazon.com/awssupport/latest/user/cost-optimization-checks.html#aws-lambda-functions-with-high-error-rates"
)

type HighErrorRateLambdaFunction struct {
	Region                  string `json:"region"`
	LambdaFunctionName      string `json:"lambdaFunctionName"`
	Reason                  string `json:"reason"`
	EstimatedMonthlySavings int    `json:"estimatedMonthlySavings"`
}

type HighErrorRateLambdaFunctionsCheck struct {
	common.Check
	HighErrorRateLambdaFunctions []HighErrorRateLambdaFunction `json:"HighErrorRateLambdaFunctions"`
}

func (v HighErrorRateLambdaFunctionsCheck) List() *HighErrorRateLambdaFunctionsCheck {
	check := &HighErrorRateLambdaFunctionsCheck{
		Check: common.Check{
			Id:                  HighErrorRateLambdaFunctionsCheckId,
			Name:                HighErrorRateLambdaFunctionsCheckName,
			Description:         HighErrorRateLambdaFunctionsCheckDescription,
			Criteria:            HighErrorRateLambdaFunctionsCheckCriteria,
			RecommendedAction:   HighErrorRateLambdaFunctionsCheckRecommendedAction,
			AdditionalResources: HighErrorRateLambdaFunctionsCheckAdditionalResources,
		},
	}
	return check
}

func (v HighErrorRateLambdaFunctionsCheck) Run(ctx context.Context, conn client.AWSClient) (*HighErrorRateLambdaFunctionsCheck, error) {
	check := new(HighErrorRateLambdaFunctionsCheck).List()

	currentTime := time.Now()

	var lambdaFunctions []HighErrorRateLambdaFunction

	in := &lambda.ListFunctionsInput{}

	paginator := lambda.NewListFunctionsPaginator(conn.Lambda, in, func(o *lambda.ListFunctionsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		for _, function := range output.Functions {
			invocations, err := conn.Cloudwatch.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
				MetricName: aws.String("Invocations"),
				Period:     aws.Int32(86400),
				Namespace:  aws.String("AWS/Lambda"),
				Statistics: []types.Statistic{
					types.StatisticSum,
				},
				Dimensions: []types.Dimension{
					{
						Name:  aws.String("FunctionName"),
						Value: function.FunctionName,
					},
				},
				StartTime: aws.Time(currentTime.AddDate(0, 0, -7)),
				EndTime:   aws.Time(currentTime),
			})

			if err != nil {
				return nil, err
			}

			sum_invocations := 0
			for _, datapoint := range invocations.Datapoints {
				sum_invocations += int(*datapoint.Sum)
			}

			if sum_invocations > 0 {
				errors, err := conn.Cloudwatch.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
					MetricName: aws.String("Errors"),
					Period:     aws.Int32(86400),
					Namespace:  aws.String("AWS/Lambda"),
					Statistics: []types.Statistic{
						types.StatisticSum,
					},
					Dimensions: []types.Dimension{
						{
							Name:  aws.String("FunctionName"),
							Value: function.FunctionName,
						},
					},
					StartTime: aws.Time(currentTime.AddDate(0, 0, -7)),
					EndTime:   aws.Time(currentTime),
				})

				if err != nil {
					fmt.Println(err)
					return nil, err
				}

				var sum_errors int
				for _, datapoint := range errors.Datapoints {
					sum_errors += int(*datapoint.Sum)
				}

				rate := sum_errors / sum_invocations * 100

				if rate > 10 {
					info := HighErrorRateLambdaFunction{
						Region:                  conn.Region,
						LambdaFunctionName:      *function.FunctionName,
						Reason:                  "HighErrors",
						EstimatedMonthlySavings: 0,
					}

					lambdaFunctions = append(lambdaFunctions, info)
				}
			}
		}
	}

	if len(lambdaFunctions) == 0 {
		return nil, nil
	}

	check.HighErrorRateLambdaFunctions = lambdaFunctions
	return check, nil
}
