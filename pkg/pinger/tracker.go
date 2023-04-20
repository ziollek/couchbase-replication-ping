package pinger

import (
	"context"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
	"time"
)

type PingTracker struct {
	replPing     interfaces.Pinger
	keyGenerator interfaces.KeyGenerator
}

func (tracker *PingTracker) Ping() (interfaces.Timing, error) {
	key := tracker.keyGenerator.Generate()
	timing := model.NewTimingRecord()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pingTiming, err := tracker.replPing.Ping(ctx, key)
	timing.AddChild("ping", pingTiming)
	if err != nil {
		return timing, err
	}
	pongTiming, err := tracker.replPing.Pong(ctx, key)
	timing.AddChild("pong", pongTiming)
	if err != nil {
		return timing, err
	}
	return timing, nil
}

func (tracker *PingTracker) OneWay() (interfaces.Timing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return tracker.replPing.Ping(ctx, tracker.keyGenerator.Generate())
}

func (tracker *PingTracker) Pong() (interfaces.Timing, error) {
	key := tracker.keyGenerator.Generate()
	timing := model.NewTimingRecord()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pongTiming, err := tracker.replPing.Pong(ctx, key)
	timing.AddChild("ping", pongTiming)
	if err != nil {
		return timing, err
	}

	pingTiming, err := tracker.replPing.Ping(ctx, key)
	timing.AddChild("ping", pingTiming)
	if err != nil {
		return timing, err
	}
	return timing, nil
}

func NewPingTracker(ping interfaces.Pinger, generator interfaces.KeyGenerator) *PingTracker {
	return &PingTracker{
		replPing:     ping,
		keyGenerator: generator,
	}
}
