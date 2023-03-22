package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AwsCmd represents the aws command
var AwsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Checks related to the aws cloud.",
	Long:  `Checks related to the aws cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify a subcommand like list or check.")
	},
}

func init() {
	rootCmd.AddCommand(AwsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// AwsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// AwsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
