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
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
}

func GetLogger() *log.Logger {
	return defaultLogger
}

func FormatByTiming(no int, timing interfaces.Timing, err error, message string) {
	FormatByTimingFields(no, GetTimingFields(timing, false), timing.GetRetries(), err, message)
}

func FormatByTwoWayTiming(no int, timing interfaces.Timing, err error, message string) {
	FormatByTimingFields(no, GetTimingFields(timing, true), timing.GetRetries(), err, message)
}

func FormatByTimingFields(no int, timingFields log.Fields, retries int, err error, message string) {
	entry := defaultLogger.WithFields(timingFields).WithField("no", no)
	if err != nil {
		entry.WithError(err).Error(message)
	} else if retries > 0 {
		entry.Warn(message)
	} else {
		entry.Info(message)
	}
}

func GetTimingFields(timing interfaces.Timing, twoWay bool) log.Fields {
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
		if twoWay {
			fields["latency"] = timing.GetPhase("wait") / 2
		} else {
			fields["latency"] = timing.GetPhase("wait") + s/2
		}
	}
	fields["total"] = timing.GetDuration()
	if timing.GetRetries() > 0 {
		fields["retries"] = timing.GetRetries()
	}
	return fields
}
