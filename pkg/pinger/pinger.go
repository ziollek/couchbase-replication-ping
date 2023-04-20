package pinger

import (
	"context"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
	"time"
)

type ReplicationPing struct {
	source      string
	destination string
	generator   interfaces.Generator
	txChannel   interfaces.Channel
	rxChannel   interfaces.Channel
}

func NewReplicationPing(
	source, destination string, generator interfaces.Generator, tx, rx interfaces.Channel,
) interfaces.Pinger {
	return &ReplicationPing{
		source,
		destination,
		generator,
		tx,
		rx,
	}
}

func (ping *ReplicationPing) Ping(ctx context.Context, key string) (interfaces.Timing, error) {
	return ping.sendAndReceive(ctx, key, ping.source, ping.txChannel)
}

func (ping *ReplicationPing) Pong(ctx context.Context, key string) (interfaces.Timing, error) {
	return ping.sendAndReceive(ctx, key, ping.destination, ping.rxChannel)
}

func (ping *ReplicationPing) sendAndReceive(
	ctx context.Context, key, source string, channel interfaces.Channel,
) (interfaces.Timing, error) {
	timing := model.NewTimingRecord()
	err := channel.Send(ping.generator.Generate(key, source))
	timing.AddPhase("send")
	if err != nil {
		return timing, err
	}

	for {
		if packet, err := channel.Receive(key); err == nil {
			if packet != nil && packet.GetOrigin() == source {
				break
			}
		}
		timing.AddPhaseTry("wait")
		select {
		case <-time.After(time.Millisecond * 1):
		case <-ctx.Done():
			return timing, ctx.Err()
		}
	}
	return timing.AddPhase("receive"), nil
}

type ReplicationHalfPing struct {
	myOrigin  string
	generator interfaces.Generator
	channel   interfaces.Channel
}

func NewReplicationHalfPing(
	myOrigin string, generator interfaces.Generator, channel interfaces.Channel,
) interfaces.Pinger {
	return &ReplicationHalfPing{
		myOrigin,
		generator,
		channel,
	}
}

func (ping *ReplicationHalfPing) Ping(ctx context.Context, key string) (interfaces.Timing, error) {
	timing := model.NewTimingRecord()
	err := ping.channel.Send(ping.generator.Generate(key, ping.myOrigin))
	return timing.AddPhase("send"), err
}

func (ping *ReplicationHalfPing) Pong(ctx context.Context, key string) (interfaces.Timing, error) {
	timing := model.NewTimingRecord()
	tries := 0
	for {
		if packet, err := ping.channel.Receive(key); err == nil {
			if packet != nil && packet.GetOrigin() != ping.myOrigin {
				break
			}
		}
		timing.AddPhaseTry("wait")
		select {
		case <-time.After(time.Millisecond * 1):
			tries++
		case <-ctx.Done():
			return timing, ctx.Err()
		}
	}
	return timing.AddPhase("receive"), nil
}
