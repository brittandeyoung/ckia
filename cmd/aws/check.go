package aws

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/brittandeyoung/ckia/cmd"
	"github.com/brittandeyoung/ckia/cmd/aws/cost"
	"github.com/spf13/cobra"
)

type Checks struct {
	CostOptimization []interface{} `json:"costOptimization"`
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
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
		fmt.Println(PrettyString(string(json)))
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
