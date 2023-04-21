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
		params := utils.ProvideParams(cmd)
		pinger, err := infra.BuildPingTracker("dynamic")
		utils.HandleError("cannot build pinger: %s", err)
		pinger.WithTimeout(params.Timeout)
		logger := utils.GetLogger()

		logger.Infof("Start measuring latency: %s", params.ToString())
		for i := 1; i <= params.Repeats; i++ {
			timing, err := pinger.OneWay()
			utils.FormatByTiming(i, timing, err, "one-way")
			time.Sleep(params.Timeout)
		}
	},
}

func init() {
	rootCmd.AddCommand(onewayCmd)
}
