package pinger

import (
	"github.com/ziollek/couchbase-replication-ping/pkg/config"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"strings"
	"time"
)

type PingPacketGenerator struct {
	params *config.Generator
}

func (g *PingPacketGenerator) Generate(key, origin string) interfaces.Packet {
	return NewPingPacket(key, origin, g.params.TTL, g.params.Size)
}

func NewPingPacketGenerator(params *config.Generator) interfaces.Generator {
	return &PingPacketGenerator{
		params: params,
	}
}

type PingPacket struct {
	Key    string        `json:"key"`
	Origin string        `json:"origin"`
	TTL    time.Duration `json:"ttl"`
	Data   string        `json:"data"`
}

func (packet *PingPacket) GetKey() string {
	return packet.Key
}

func (packet *PingPacket) GetOrigin() string {
	return packet.Origin
}

func (packet *PingPacket) GetSize() int {
	return len(packet.Data)
}

func (packet *PingPacket) GetTTL() time.Duration {
	return packet.TTL
}

func NewPingPacket(key, origin string, ttl time.Duration, size int) interfaces.Packet {
	return &PingPacket{
		key,
		origin,
		ttl,
		generateString(size),
	}
}

func generateString(size int) string {
	sb := strings.Builder{}
	for i := 0; i < size; i++ {
		sb.WriteString("1")
	}
	return sb.String()
}
