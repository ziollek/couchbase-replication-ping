package pinger

import (
	"context"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
	"time"
)

const DefaultTimeout = 10 * time.Second

type PingTracker struct {
	replPing        interfaces.Pinger
	keyGenerator    interfaces.KeyGenerator
	timeout         time.Duration
	combineStrategy bool
}

func (tracker *PingTracker) Ping() (interfaces.Timing, error) {
	key := tracker.keyGenerator.Generate()
	timing := model.NewTimingRecord()
	ctx, cancel := context.WithTimeout(context.Background(), tracker.timeout)
	defer cancel()
	pingTiming, err := tracker.replPing.Ping(ctx, key)
	tracker.mergeTimings("ping", timing, pingTiming)
	if err != nil {
		return timing, err
	}
	pongTiming, err := tracker.replPing.Pong(ctx, key)
	tracker.mergeTimings("pong", timing, pongTiming)
	if err != nil {
		return timing, err
	}
	return timing, nil
}

func (tracker *PingTracker) OneWay() (interfaces.Timing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tracker.timeout)
	defer cancel()
	return tracker.replPing.Ping(ctx, tracker.keyGenerator.Generate())
}

func (tracker *PingTracker) Pong() (interfaces.Timing, error) {
	key := tracker.keyGenerator.Generate()
	timing := model.NewTimingRecord()
	ctx, cancel := context.WithTimeout(context.Background(), tracker.timeout)
	defer cancel()
	pongTiming, err := tracker.replPing.Pong(ctx, key)
	tracker.mergeTimings("pong", timing, pongTiming)
	if err != nil {
		return timing, err
	}

	pingTiming, err := tracker.replPing.Ping(ctx, key)
	tracker.mergeTimings("ping", timing, pingTiming)
	if err != nil {
		return timing, err
	}
	return timing, nil
}

func (tracker *PingTracker) mergeTimings(phase string, parent, child interfaces.Timing) interfaces.Timing {
	if tracker.combineStrategy {
		parent.Combine(child)
	} else {
		parent.AddChild(phase, child)
	}
	return parent
}

func (tracker *PingTracker) WithTimeout(timeout time.Duration) {
	tracker.timeout = timeout
}

func (tracker *PingTracker) WithCombineTimings(combineStrategy bool) {
	tracker.combineStrategy = combineStrategy
}

func NewPingTracker(ping interfaces.Pinger, generator interfaces.KeyGenerator) *PingTracker {
	return &PingTracker{
		replPing:        ping,
		keyGenerator:    generator,
		timeout:         DefaultTimeout,
		combineStrategy: false,
	}
}
