package aws

import (
	"encoding/json"
	"fmt"

	"github.com/brittandeyoung/ckia/cmd"
	"github.com/brittandeyoung/ckia/cmd/aws/cost"
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
		allChecks := Checks{}
		allChecks.CostOptimization = append(allChecks.CostOptimization, cost.ListRDSIdleDB())
		json, err := json.Marshal(allChecks)
		if err != nil {
			fmt.Print("An Error happened when marshaling json")
		}
		fmt.Println(common.PrettyString(string(json)))
	},
}

func init() {
	cmd.AwsCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
