package pinger_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces/mocks"
	"testing"
)

func TestTrackerShouldCallPingAndPongConsecutively(t *testing.T) {
	key := "test"
	// Given: ping tracker built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	replicationPing := mocks.NewMockPinger(ctrl)
	replicationPing.EXPECT().Ping(gomock.Any(), key)
	replicationPing.EXPECT().Pong(gomock.Any(), key)

	tracker := pinger.NewPingTracker(key, replicationPing)
	// When: calling ping
	_, err := tracker.Ping()

	// Then: Ping & Pong should be called internally
	require.NoError(t, err)
}
