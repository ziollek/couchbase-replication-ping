package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ziollek/couchbase-replication-ping/internal/cmd/utils"
	"github.com/ziollek/couchbase-replication-ping/pkg/infra"
	"time"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "measure two way replication latency",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		params := utils.ProvideParams(cmd)
		pinger, err := infra.BuildPingTracker("static")
		utils.HandleError("cannot build pinger: %s", err)
		pinger.WithTimeout(params.Timeout)
		logger := utils.GetLogger()

		logger.Infof("Start measuring latency: %s", params.ToString())
		for i := 1; i <= params.Repeats; i++ {
			timing, err := pinger.Ping()
			utils.FormatByTiming(i, timing, err, "ping")
			time.Sleep(params.Interval)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
