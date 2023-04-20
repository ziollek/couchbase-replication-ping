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
		pinger, err := infra.BuildPingTracker("static")
		utils.HandleError("cannot build pinger: %s", err)
		n, err := cmd.Flags().GetInt("repeat")
		utils.HandleError("improper repeat option: %s", err)
		interval, err := cmd.Flags().GetDuration("interval")
		utils.HandleError("improper interval option: %s", err)
		logger := utils.GetLogger()

		logger.Info("Start measuring latency ... ")
		for i := 1; i <= n; i++ {
			timing, err := pinger.Ping()
			utils.FormatByTiming(i, timing, err, "ping")
			time.Sleep(interval)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.PersistentFlags().Int("repeat", 3, "define how many times ping should be repeated")
	pingCmd.PersistentFlags().Duration("interval", time.Second, "define pings frequency, default every 1s")
}
