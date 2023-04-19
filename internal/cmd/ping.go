package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ziollek/couchbase-replication-ping/pkg/config"
	"github.com/ziollek/couchbase-replication-ping/pkg/infra"
	"os"
	"time"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "measure two way replication latency",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.GetConfig()
		handleError("cannot load config: %s", err)
		pinger, err := infra.Build(c)
		handleError("cannot build pinger: %s", err)
		n, err := cmd.Flags().GetInt("repeat")
		handleError("improper repeat option: %s", err)

		fmt.Println("Start measuring latency ... ")
		for i := 1; i <= n; i++ {
			duration, err := pinger.Ping()
			fmt.Printf("[%s] ping %d)\tduration: %s, err: %s\n", time.Now().Format("2006-01-02 15:04:05"), i, duration, err)
			time.Sleep(time.Second)
		}
	},
}

func handleError(format string, err error) {
	if err != nil {
		fatal(fmt.Sprintf(format, err))
	}
}

func fatal(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.PersistentFlags().Int("repeat", 3, "how many times ping should be repeated")
}
