package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run available checks for aws",
	Long:  `Run available opinionated checks for aws cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		conn := client.InitiateClient(cfg)

		allChecks := Checks{}
		checksMap := internalAws.BuildChecksMap()
		for k := range checksMap {
			if strings.Contains(k, "aws:cost") {
				res, err := common.Call(k, checksMap, common.MethodNameRun, ctx, conn)
				if err != nil {
					log.Fatal(err)
				}
				if res != nil {
					allChecks.CostOptimization = append(allChecks.CostOptimization, res)
				}
			}
			if strings.Contains(k, "aws:security") {
				res, err := common.Call(k, checksMap, common.MethodNameRun, ctx, conn)
				if err != nil {
					log.Fatal(err)
				}
				if res != nil {
					allChecks.Security = append(allChecks.Security, res)
				}
			}
		}

		json, err := json.Marshal(allChecks)
		if err != nil {
			fmt.Print("An Error happened when marshaling json")
		}
		resp, err := common.PrettyString(string(json))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}

func init() {
	cmd.AwsCmd.AddCommand(checkCmd)
}
