package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/brittandeyoung/ckia/cmd"
	internalAws "github.com/brittandeyoung/ckia/internal/aws"
	"github.com/brittandeyoung/ckia/internal/client"
	"github.com/brittandeyoung/ckia/internal/common"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/smirzaei/parallel"
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
var outFile string
var outFormat string

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run available checks for aws",
	Long:  `Run available opinionated checks for aws cloud.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags
		if !common.StringSliceContains([]string{"json"}, outFormat) {
			return errors.New("unsupported format provided to out-format flag")
		}

		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		conn := client.InitiateClient(cfg)
		allChecks := Checks{}
		checksMap := internalAws.BuildChecksMap()
		checksList := []string{}
		for k := range checksMap {
			checksList = append(checksList, k)
		}
		bar := progressbar.NewOptions(len(checksList),
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionFullWidth(),
			progressbar.OptionShowCount(),
			progressbar.OptionShowElapsedTimeOnFinish(),
			progressbar.OptionSetDescription(fmt.Sprintf("Running [cyan][%d][reset] ckia Checks...", len(checksList))),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}))
		errors := parallel.Map(checksList, func(k string) error {
			if (len(includeChecks) > 0 && common.StringSliceContains(includeChecks, k)) || (len(excludeChecks) > 0 && !common.StringSliceContains(excludeChecks, k)) || (len(includeChecks) == 0 && len(excludeChecks) == 0) {
				if strings.Contains(k, "aws:cost") {
					res, _ := common.Call(k, checksMap, common.MethodNameRun, ctx, conn)
					if err != nil {
						return err
					}
					if res != nil {
						allChecks.CostOptimization = append(allChecks.CostOptimization, res)
					}
				}
				if strings.Contains(k, "aws:security") {
					res, _ := common.Call(k, checksMap, common.MethodNameRun, ctx, conn)
					if err != nil {
						return err
					}
					if res != nil {
						allChecks.Security = append(allChecks.Security, res)
					}
				}
			}
			bar.Add(1)
			return nil
		})

		for _, err := range errors {
			if err != nil {
				return err
			}
		}

		json, err := json.Marshal(allChecks)
		if err != nil {
			fmt.Print("An Error happened when marshaling json")
		}
		fmt.Println()
		resp, err := common.PrettyString(string(json))
		if err != nil {
			return err
		}

		if outFile != "" {
			err = ioutil.WriteFile(outFile, json, 0644)

			if err != nil {
				return err
			}
		} else {
			fmt.Println(resp)
		}

		return nil
	},
}

func init() {
	cmd.AwsCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringSliceVarP(&includeChecks, "include-checks", "i", []string{}, "A list of all the checks you wish to run.")
	checkCmd.Flags().StringSliceVarP(&excludeChecks, "exclude-checks", "e", []string{}, "A list of checks to exclude from running.")
	checkCmd.Flags().StringVarP(&outFile, "out-file", "o", "", "A path to a file to store check results.")
	checkCmd.Flags().StringVarP(&outFormat, "out-format", "f", "json", "The file output format for check results. Default: json (currently only json is supported).")
}
