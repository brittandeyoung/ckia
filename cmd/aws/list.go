package aws

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/brittandeyoung/ckia/cmd"
	internalAws "github.com/brittandeyoung/ckia/internal/aws"
	"github.com/brittandeyoung/ckia/internal/common"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available checks for aws",
	Long:  `List the available opinionated checks for aws cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		allChecks := Checks{}
		checksMap := internalAws.BuildChecksMap()
		for k := range checksMap {
			if strings.Contains(k, "aws:cost") {
				res, _ := common.Call(k, checksMap, common.MethodNameList)
				if res != nil {
					allChecks.CostOptimization = append(allChecks.CostOptimization, res)
				}
			}
			if strings.Contains(k, "aws:security") {
				res, _ := common.Call(k, checksMap, common.MethodNameList)
				if res != nil {
					allChecks.Security = append(allChecks.Security, res)
				}
			}
		}
		json, err := json.Marshal(allChecks)
		if err != nil {
			fmt.Print("An Error happened when marshaling json")
		}
		fmt.Println(common.PrettyString(string(json)))
	},
}

func init() {
	cmd.AwsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
