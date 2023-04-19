package interfaces

import (
	"context"
	"time"
)

type Packet interface {
	GetKey() string
	GetOrigin() string
	GetTTL() time.Duration
	GetSize() int
}

type Generator interface {
	Generate(string, string) Packet
}

type Channel interface {
	Send(packet Packet) error
	Receive(key string) (Packet, error)
}

type Pinger interface {
	Ping(context.Context, string) error
	Pong(context.Context, string) error
}
