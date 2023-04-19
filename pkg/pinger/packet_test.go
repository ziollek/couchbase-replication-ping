package pinger_test

import (
	"github.com/stretchr/testify/require"
	"github.com/ziollek/couchbase-replication-ping/pkg/config"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger"
	"testing"
	"time"
)

func TestShouldGeneratePacketWithRequestedLength(t *testing.T) {
	size := 10
	// Given: Packet that satisfies interface
	packet := pinger.NewPingPacket("test", "origin", time.Minute, size)
	require.Equal(t, size, packet.GetSize())

	// When: cast it to concrete struct
	casted, ok := packet.(*pinger.PingPacket)
	require.True(t, ok)

	// Then: data should have expected length
	require.Equal(t, size, len(casted.Data))
}

func TestGeneratorShouldGeneratePacketWithApropriateParams(t *testing.T) {
	// Given: configuration & generator based on it
	generatorConfig := &config.Generator{
		TTL:  time.Second,
		Size: 20,
	}
	generator := pinger.NewPingPacketGenerator(generatorConfig)

	// When: generating new packet
	packet := generator.Generate("key", "origin")

	// Then: it should be composed from proper values
	require.Equal(t, "key", packet.GetKey())
	require.Equal(t, "origin", packet.GetOrigin())
	require.Equal(t, generatorConfig.TTL, packet.GetTTL())
	require.Equal(t, generatorConfig.Size, packet.GetSize())
}
