package model_test

import (
	"github.com/stretchr/testify/require"
	"github.com/ziollek/couchbase-replication-ping/pkg/config"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
	"strings"
	"testing"
	"time"
)

func TestShouldGenerateStaticKey(t *testing.T) {
	key := "key"
	// Given: static generator
	generator := model.NewStaticKeyGenerator(key)

	// When: generating key
	result := generator.Generate()

	// Then: generate key should the same as passed to constructor
	require.Equal(t, key, result)
}

func TestShouldGenerateRandomKeys(t *testing.T) {
	prefix := "prefix"
	length := 4

	// Given: configured random generator
	generator := model.NewRandomKeyGenerator(prefix, length)

	// When: generating two consecutive keys
	first := generator.Generate()
	second := generator.Generate()

	// Then: generate keys should be different
	require.NotEqual(t, first, second)
	// Then: keys should start with prefix
	require.True(t, strings.HasPrefix(first, prefix))
	require.True(t, strings.HasPrefix(second, prefix))
	// Then: and should have specified length plus one for delimiter
	require.Len(t, first, len(prefix)+length+1)
	require.Len(t, second, len(prefix)+length+1)
}

func TestShouldGeneratePacketWithRequestedLength(t *testing.T) {
	size := 10
	// Given: Packet that satisfies interface
	packet := model.NewPingPacket("test", "origin", time.Minute, size)
	require.Equal(t, size, packet.GetSize())

	// When: cast it to concrete struct
	casted, ok := packet.(*model.PingPacket)
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
	generator := model.NewPingPacketGenerator(generatorConfig)

	// When: generating new packet
	packet := generator.Generate("key", "origin")

	// Then: it should be composed from proper values
	require.Equal(t, "key", packet.GetKey())
	require.Equal(t, "origin", packet.GetOrigin())
	require.Equal(t, generatorConfig.TTL, packet.GetTTL())
	require.Equal(t, generatorConfig.Size, packet.GetSize())
}
