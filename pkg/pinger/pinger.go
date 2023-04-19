package pinger

import (
	"context"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
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

func (ping *ReplicationPing) Ping(ctx context.Context, key string) error {
	if err := ping.txChannel.Send(ping.generator.Generate(key, ping.source)); err != nil {
		return err
	}
	for {
		if packet, err := ping.rxChannel.Receive(key); err == nil {
			if packet != nil && packet.GetOrigin() == ping.source {
				break
			}
		}
		select {
		case <-time.After(time.Millisecond * 1):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (ping *ReplicationPing) Pong(ctx context.Context, key string) error {
	if err := ping.rxChannel.Send(ping.generator.Generate(key, ping.destination)); err != nil {
		return err
	}
	for {
		if packet, err := ping.txChannel.Receive(key); err == nil {
			if packet != nil && packet.GetOrigin() == ping.destination {
				break
			}
		}
		select {
		case <-time.After(time.Millisecond * 1):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
