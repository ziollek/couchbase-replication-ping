package cmd

import (
	"github.com/ziollek/couchbase-replication-ping/internal/cmd/utils"
	"github.com/ziollek/couchbase-replication-ping/pkg/infra"
	"time"

	"github.com/spf13/cobra"
)

// onewayCmd represents the oneway command
var onewayCmd = &cobra.Command{
	Use:   "oneway",
	Short: "measure one way replication latency",
	Long:  `It gives additional glimpse about write/read times from source and destination bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		pinger, err := infra.BuildPingTracker("dynamic")
		utils.HandleError("cannot build pinger: %s", err)
		n, err := cmd.Flags().GetInt("repeat")
		utils.HandleError("improper repeat option: %s", err)
		interval, err := cmd.Flags().GetDuration("interval")
		utils.HandleError("improper interval option: %s", err)
		logger := utils.GetLogger()

		logger.Info("Start measuring latency ... ")
		for i := 1; i <= n; i++ {
			timing, err := pinger.OneWay()
			utils.FormatByTiming(i, timing, err, "one-way")
			time.Sleep(interval)
		}
	},
}

func init() {
	rootCmd.AddCommand(onewayCmd)
	onewayCmd.PersistentFlags().Int("repeat", 3, "define how many times ping should be repeated")
	onewayCmd.PersistentFlags().Duration("interval", time.Second, "define pings frequency, default every 1s")
}
