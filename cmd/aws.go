/*
Copyright © 2023 Brittan DeYoung <brittan@deyoung.dev>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
		fmt.Println("Error: must also specify a subcommand like list, run.")
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
