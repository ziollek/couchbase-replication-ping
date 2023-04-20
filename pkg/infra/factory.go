package infra

import (
	"github.com/ziollek/couchbase-replication-ping/pkg/config"
	"github.com/ziollek/couchbase-replication-ping/pkg/kv"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
)

func BuildPingTracker(keyMode string) (*pinger.PingTracker, error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	source, err := kv.NewKV(c.Source)
	if err != nil {
		return nil, err
	}
	destination, _ := kv.NewKV(c.Destination)
	if err != nil {
		return nil, err
	}
	replPing := pinger.NewReplicationPing(
		c.Source.Name,
		c.Destination.Name,
		model.NewPingPacketGenerator(c.Generator),
		pinger.NewChannel(source, destination),
		pinger.NewChannel(destination, source),
	)
	return pinger.NewPingTracker(replPing, getKeyGeneratorByMode(c.Key, keyMode)), nil
}

func BuildHalfPingTracker(origin string) (*pinger.PingTracker, error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	originKV, err := kv.NewKV(getCouchbaseByOrigin(origin, c))
	if err != nil {
		return nil, err
	}
	replPing := pinger.NewReplicationHalfPing(
		origin,
		model.NewPingPacketGenerator(c.Generator),
		pinger.NewChannel(originKV, originKV),
	)
	return pinger.NewPingTracker(replPing, getKeyGeneratorByMode(c.Key, "static")), nil
}

func getCouchbaseByOrigin(origin string, c *config.Config) *config.Couchbase {
	if origin == "source" {
		return c.Source
	}
	return c.Destination
}

func getKeyGeneratorByMode(key, mode string) interfaces.KeyGenerator {
	if mode == "static" {
		return model.NewStaticKeyGenerator(key)
	}
	return model.NewRandomKeyGenerator(key, 5)
}
