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
	RunE: func(cmd *cobra.Command, args []string) error {
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
		resp, err := common.PrettyString(string(json))
		if err != nil {
			return err
		}
		fmt.Println(resp)
		return nil
	},
}

func init() {
	cmd.AwsCmd.AddCommand(listCmd)
}
