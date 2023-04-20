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

		pinger, err := infra.BuildHalfPingTracker(origin)
		utils.HandleError("cannot build pinger: %s", err)
		n, err := cmd.Flags().GetInt("repeat")
		utils.HandleError("improper repeat option: %s", err)
		interval, err := cmd.Flags().GetDuration("interval")
		utils.HandleError("improper interval option: %s", err)

		logger.Infof("%v", args)

		logger.Infof("Start measuring latency from %s perspective ... ", origin)
		for i := 1; i <= n; i++ {
			var timing interfaces.Timing
			if origin == "source" {
				timing, err = pinger.Ping()
			} else {
				timing, err = pinger.Pong()
			}
			utils.FormatByTiming(i, timing, err, origin)
			time.Sleep(interval)
		}
	},
}

func init() {
	rootCmd.AddCommand(halfpingCmd)

	halfpingCmd.PersistentFlags().Int("repeat", 3, "define how many times ping should be repeated")
	halfpingCmd.PersistentFlags().Duration("interval", time.Second, "define pings frequency, default every 1s")
}
