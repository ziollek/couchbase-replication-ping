package pinger

import (
	kv "github.com/ziollek/couchbase-replication-ping/pkg/kv/interfaces"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
)

type KVChannel struct {
	source      kv.KV
	destination kv.KV
}

func (channel *KVChannel) Send(packet interfaces.Packet) error {
	return channel.source.Upsert(packet.GetKey(), packet, packet.GetTTL())
}

func (channel *KVChannel) Receive(key string) (interfaces.Packet, error) {
	var p model.PingPacket
	err := channel.destination.Get(key, &p)
	return &p, err
}

func NewChannel(source, destination kv.KV) interfaces.Channel {
	return &KVChannel{
		source:      source,
		destination: destination,
	}
}
