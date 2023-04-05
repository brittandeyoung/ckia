package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/brittandeyoung/ckia/cmd"
	internalAws "github.com/brittandeyoung/ckia/internal/aws"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/common"
	"github.com/spf13/cobra"
)

type Checks struct {
	CostOptimization []interface{} `json:"costOptimization"`
	Performance      []interface{} `json:"performance"`
	Security         []interface{} `json:"security"`
	FaultTolerance   []interface{} `json:"faultTolerance"`
	ServiceLimits    []interface{} `json:"serviceLimits"`
}

var includeChecks []string
var excludeChecks []string

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run available checks for aws",
	Long:  `Run available opinionated checks for aws cloud.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		conn := client.InitiateClient(cfg)
		allChecks := Checks{}
		checksMap := internalAws.BuildChecksMap()
		for k := range checksMap {
			if (len(includeChecks) > 0 && common.StringSliceContains(includeChecks, k)) || (len(excludeChecks) > 0 && !common.StringSliceContains(excludeChecks, k)) {
				if strings.Contains(k, "aws:cost") {
					res, err := common.Call(k, checksMap, common.MethodNameRun, ctx, conn)
					if err != nil {
						return err
					}
					if res != nil {
						allChecks.CostOptimization = append(allChecks.CostOptimization, res)
					}
				}
				if strings.Contains(k, "aws:security") {
					res, err := common.Call(k, checksMap, common.MethodNameRun, ctx, conn)
					if err != nil {
						return err
					}
					if res != nil {
						allChecks.Security = append(allChecks.Security, res)
					}
				}
			}
		}

		json, err := json.Marshal(allChecks)
		if err != nil {
			fmt.Print("An Error happened when marshaling json")
		}
		resp, err := common.PrettyString(string(json))
		if err != nil {
			return err
		}
		fmt.Println(resp)
		return nil
	},
}

func init() {
	cmd.AwsCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringSliceVarP(&includeChecks, "include-checks", "i", []string{}, "A list of all the checks you wish to run.")
	checkCmd.Flags().StringSliceVarP(&excludeChecks, "exclude-checks", "e", []string{}, "A list of checks to exclude from running.")
}
