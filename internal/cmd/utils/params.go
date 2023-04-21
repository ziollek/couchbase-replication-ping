package utils

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

type GenericParams struct {
	Interval time.Duration
	Repeats  int
	Timeout  time.Duration
}

func (gp *GenericParams) ToString() string {
	return fmt.Sprintf("interval=%s, repeats=%d, timeout=%s", gp.Interval, gp.Repeats, gp.Timeout)
}

func ProvideParams(cmd *cobra.Command) *GenericParams {
	interval, err := cmd.Flags().GetDuration("interval")
	HandleError("improper interval option: %s", err)
	n, err := cmd.Flags().GetInt("repeat")
	HandleError("improper repeat option: %s", err)
	timeout, err := cmd.Flags().GetDuration("timeout")
	HandleError("improper timeout option: %s", err)
	return &GenericParams{
		Repeats:  n,
		Timeout:  timeout,
		Interval: interval,
	}
}
