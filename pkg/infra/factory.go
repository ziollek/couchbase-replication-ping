package infra

import (
	"github.com/ziollek/couchbase-replication-ping/pkg/config"
	"github.com/ziollek/couchbase-replication-ping/pkg/kv"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger"
)

func Build(config *config.Config) (*pinger.PingTracker, error) {
	source, err := kv.NewKV(config.Source)
	if err != nil {
		return nil, err
	}
	destination, _ := kv.NewKV(config.Source)
	if err != nil {
		return nil, err
	}
	replPing := pinger.NewReplicationPing(
		config.Source.Name,
		config.Destination.Name,
		pinger.NewPingPacketGenerator(config.Generator),
		pinger.NewChannel(source, destination),
		pinger.NewChannel(destination, source),
	)
	return pinger.NewPingTracker(config.Key, replPing), nil
}
