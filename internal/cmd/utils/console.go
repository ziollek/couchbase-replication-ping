package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"os"
	"time"
)

var defaultLogger *log.Logger

func ConfigureLogger(json bool) {
	defaultLogger = log.New()
	defaultLogger.SetOutput(os.Stdout)
	defaultLogger.SetLevel(log.InfoLevel)
	if json {
		defaultLogger.SetFormatter(&log.JSONFormatter{})
	} else {
		defaultLogger.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}
}

func GetLogger() *log.Logger {
	return defaultLogger
}

func FormatByTiming(no int, timing interfaces.Timing, err error, message string) {
	entry := defaultLogger.WithFields(getTimingFields(timing)).WithField("no", no)
	if err != nil {
		entry.WithError(err).Error(message)
	} else if timing.GetRetries() > 0 {
		entry.Warn(message)
	} else {
		entry.Info(message)
	}
}

func getTimingFields(timing interfaces.Timing) log.Fields {
	fields := log.Fields{}
	s := time.Duration(0)
	n := 0
	for _, phase := range timing.GetPhases() {
		fields[phase] = timing.GetPhase(phase)
		if phase != "wait" {
			s += timing.GetPhase(phase)
			n++
		}
	}
	if timing.GetPhase("wait") > 0 && n == 2 {
		fields["latency"] = timing.GetPhase("wait") + s/2
	}
	fields["total"] = timing.GetDuration()
	if timing.GetRetries() > 0 {
		fields["retries"] = timing.GetRetries()
	}
	return fields
}
