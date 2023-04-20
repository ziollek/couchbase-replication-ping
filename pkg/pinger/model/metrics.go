package model

import (
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces"
	"golang.org/x/exp/maps"
	"time"
)

// TimingRecord represents detailed timing of multistage/multiphase operation
type TimingRecord struct {
	offset  time.Time
	phases  map[string]time.Duration
	retries map[string]int
}

func (record *TimingRecord) AddPhase(phase string) interfaces.Timing {
	record.phases[phase] = time.Since(record.offset)
	record.offset = time.Now()
	return record
}

func (record *TimingRecord) AddPhaseTry(phase string) interfaces.Timing {
	if _, ok := record.phases[phase]; !ok {
		record.phases[phase] = time.Duration(0)
		record.retries[phase] = 0
	}
	record.phases[phase] += time.Since(record.offset)
	record.retries[phase]++
	record.offset = time.Now()
	return record
}

func (record *TimingRecord) AddChild(phase string, timing interfaces.Timing) {
	record.phases[phase] = timing.GetDuration()
	record.retries[phase] = timing.GetRetries()
}

func (record *TimingRecord) GetPhases() []string {
	return maps.Keys(record.phases)
}

func (record *TimingRecord) GetPhase(phase string) time.Duration {
	if result, ok := record.phases[phase]; ok {
		return result
	}
	return 0
}

func (record *TimingRecord) GetDuration() (result time.Duration) {
	for _, duration := range record.phases {
		result += duration
	}
	return result
}

func (record *TimingRecord) GetRetries() (result int) {
	for _, retries := range record.retries {
		result += retries
	}
	return result
}

func NewTimingRecord() interfaces.Timing {
	return &TimingRecord{
		time.Now(),
		map[string]time.Duration{},
		map[string]int{},
	}
}
