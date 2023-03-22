package aws

import (
	"fmt"

	"github.com/brittandeyoung/ckia/cmd"
	internalAws "github.com/brittandeyoung/ckia/internal/aws"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available checks for aws",
	Long:  `List the available opinionated checks for aws cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		checksMap := internalAws.BuildChecksMap()
		for k, _ := range checksMap {
			fmt.Println(k)
		}
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
