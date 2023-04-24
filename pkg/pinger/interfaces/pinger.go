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

type Timing interface {
	AddPhase(phase string) Timing
	AddPhaseTry(phase string) Timing
	AddChild(phase string, timing Timing)
	Combine(timing Timing)
	GetPhases() []string
	GetPhase(phase string) time.Duration
	GetPhaseRetries(phase string) int
	GetDuration() time.Duration
	GetRetries() int
}

type Generator interface {
	Generate(key, origin string) Packet
}

type KeyGenerator interface {
	Generate() string
}

type Channel interface {
	Send(packet Packet) error
	Receive(key string) (Packet, error)
}

type Pinger interface {
	Ping(ctx context.Context, key string) (Timing, error)
	Pong(ctx context.Context, key string) (Timing, error)
}
