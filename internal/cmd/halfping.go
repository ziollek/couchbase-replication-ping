/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ziollek/couchbase-replication-ping/internal/cmd/utils"
	"github.com/ziollek/couchbase-replication-ping/pkg/infra"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"time"

	"github.com/spf13/cobra"
)

// halfpingCmd represents the halfping command
var halfpingCmd = &cobra.Command{
	Use:   "halfping",
	Short: "make ping requests and expects pong from another side",
	Long: `This mode required running programs in two independent terminals:
- with mode source
- with mode destination
The source mode variant stores the document in the source bucket, the destination mode variant tries to read document from the destination bucket.
When a document is read by the destination variant it stores them back with the changed value in order to source mode variant can detect change
This logic is repeated as many times as is defined by the repeat flag`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"source", "destination"},
	Run: func(cmd *cobra.Command, args []string) {
		logger := utils.GetLogger()
		if len(args) != 1 {
			logger.Fatal("Exactly one argument was expected")
		}
		origin := args[0]
		logger.Infof("Start measuring latency from %s perspective ... ", origin)
		params := utils.ProvideParams(cmd)
		pinger, err := infra.BuildHalfPingTracker(origin)
		utils.HandleError("cannot build pinger: %s", err)
		pinger.WithTimeout(params.Timeout)

		logger.Infof("Start measuring latency from %s perspective: %s", origin, params.ToString())
		for i := 1; i <= params.Repeats; i++ {
			var timing interfaces.Timing
			if origin == "source" {
				timing, err = pinger.Ping()
			} else {
				timing, err = pinger.Pong()
			}
			utils.FormatByTiming(i, timing, err, origin)
			time.Sleep(params.Interval)
		}
	},
}

func init() {
	rootCmd.AddCommand(halfpingCmd)
}
