package pinger

import (
	"context"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"time"
)

type PingTracker struct {
	key      string
	replPing interfaces.Pinger
}

func (tracker *PingTracker) Ping() (time.Duration, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := tracker.replPing.Ping(ctx, tracker.key); err != nil {
		return time.Since(start), err
	}
	if err := tracker.replPing.Pong(ctx, tracker.key); err != nil {
		return time.Since(start), err
	}
	return time.Since(start), nil
}

func NewPingTracker(key string, ping interfaces.Pinger) *PingTracker {
	return &PingTracker{
		key:      key,
		replPing: ping,
	}
}
