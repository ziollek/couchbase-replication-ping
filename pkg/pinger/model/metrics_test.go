package model_test

import (
	"github.com/stretchr/testify/require"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
	"testing"
	"time"
)

func TestShouldSumDurationForAllPhases(t *testing.T) {
	// Given: timing record
	record := model.NewTimingRecord()

	// When: adding phases
	record.AddPhase("one")
	record.AddPhase("two")

	// Then: all phases should be sum up
	require.ElementsMatch(t, []string{"one", "two"}, record.GetPhases())
	require.True(t, record.GetDuration() > 0)
	// Then: and there could be ability to fetch duration for single phase
	require.True(t, record.GetPhase("two") > 0)
	// Then: and there should be duration=0 returned for unknown phase
	require.Equal(t, time.Duration(0), record.GetPhase("three"))
}

func TestShouldSumsRetries(t *testing.T) {
	// Given: timing record
	record := model.NewTimingRecord()
	// When: retrying phases
	record.AddPhaseTry("one")
	record.AddPhaseTry("one")
	record.AddPhaseTry("two")
	// Then: all retried phases should be sum up
	require.ElementsMatch(t, []string{"one", "two"}, record.GetPhases())
	require.True(t, record.GetDuration() > 0)
	// Then: and sum of all retries should be returned
	require.Equal(t, 3, record.GetRetries())
}

func TestShouldIncludeChildGeneralMetricsAndSumRetries(t *testing.T) {
	// Given: timing records parent & child
	parent := model.NewTimingRecord()
	child := model.NewTimingRecord()

	// When: retrying phases on both
	parent.AddPhaseTry("general")
	child.AddPhaseTry("one")
	child.AddPhaseTry("two")
	// And: add child to parent
	parent.AddChild("specific", child)

	// Then: parent phases and one that represents child should be available
	require.ElementsMatch(t, []string{"general", "specific"}, parent.GetPhases())
	require.True(t, parent.GetDuration() > 0)
	// Then: and sum of all retries (parent & child) should be returned
	require.Equal(t, 3, parent.GetRetries())
}
