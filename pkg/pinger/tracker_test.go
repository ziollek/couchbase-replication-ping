package pinger_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces/mocks"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
	"testing"
)

func TestTrackerShouldCallPingAndPongConsecutivelyOnPingCall(t *testing.T) {
	key := "test"
	// Given: ping tracker built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	keyGenerator := mocks.NewMockKeyGenerator(ctrl)
	keyGenerator.EXPECT().Generate().Return(key)

	replicationPing := mocks.NewMockPinger(ctrl)
	replicationPing.EXPECT().Ping(gomock.Any(), key).Return(&model.TimingRecord{}, nil)
	replicationPing.EXPECT().Pong(gomock.Any(), key).Return(&model.TimingRecord{}, nil)

	tracker := pinger.NewPingTracker(replicationPing, keyGenerator)
	// When: calling ping
	_, err := tracker.Ping()

	// Then: Ping & Pong should be called internally
	require.NoError(t, err)
}

func TestTrackerShouldCallPongAndPingConsecutivelyOnPongCall(t *testing.T) {
	key := "test"
	// Given: ping tracker built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	keyGenerator := mocks.NewMockKeyGenerator(ctrl)
	keyGenerator.EXPECT().Generate().Return(key)

	replicationPing := mocks.NewMockPinger(ctrl)
	replicationPing.EXPECT().Pong(gomock.Any(), key).Return(&model.TimingRecord{}, nil)
	replicationPing.EXPECT().Ping(gomock.Any(), key).Return(&model.TimingRecord{}, nil)

	tracker := pinger.NewPingTracker(replicationPing, keyGenerator)
	// When: calling pong
	_, err := tracker.Pong()

	// Then: Ping & Pong should be called internally
	require.NoError(t, err)
}

func TestTrackerShouldCallPingOnOneWay(t *testing.T) {
	key := "test"
	timing := &model.TimingRecord{}
	// Given: ping tracker built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	keyGenerator := mocks.NewMockKeyGenerator(ctrl)
	keyGenerator.EXPECT().Generate().Return(key)

	replicationPing := mocks.NewMockPinger(ctrl)
	replicationPing.EXPECT().Ping(gomock.Any(), key).Return(timing, nil)

	tracker := pinger.NewPingTracker(replicationPing, keyGenerator)
	// When: calling oneway
	result, err := tracker.OneWay()

	// Then: only Ping should be called internally
	require.NoError(t, err)
	require.Equal(t, timing, result)
}
