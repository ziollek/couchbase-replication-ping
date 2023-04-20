package cmd

import (
	"github.com/ziollek/couchbase-replication-ping/internal/cmd/utils"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	jsonOutput bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cb-tracker",
	Short: "cb-tracker is tool for measuring replication latencies in variety of ways",
	Long: `It operates analogically as ping tool for examining network round-trip-time'
In such a case:
- network is set of couchbase clusters
- host is single cluster
- link is XDCR that transfers data from one host to another

By connecting to both side of connection, the tool allows to measure
how much time it takes from storing document in source bucket 
till receiving it from destination one and vice-versa.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cb-tracker.yaml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output as a json")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	utils.ConfigureLogger(jsonOutput)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".internal" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cb-tracker")
	}

	viper.AutomaticEnv() // read in environment variables that match
	logger := utils.GetLogger()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Infof("Using config file: %s", viper.ConfigFileUsed())
	} else {
		logger.Fatal("Cannot find any config")
	}
}
