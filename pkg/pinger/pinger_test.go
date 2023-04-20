package pinger_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces/mocks"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/model"
	"testing"
	"time"
)

func TestReplicationPingShouldSendGeneratedPacketViaTxAndReceiveItViaTxOnPing(t *testing.T) {
	// Given: replication ping built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	packet := model.NewPingPacket("test", "source", time.Second, 1)
	generatorMock := mocks.NewMockGenerator(ctrl)
	generatorMock.EXPECT().Generate("test", "source").Return(packet)
	txMock := mocks.NewMockChannel(ctrl)
	txMock.EXPECT().Send(packet).Return(nil)
	txMock.EXPECT().Receive("test").Return(packet, nil)
	replPing := pinger.NewReplicationPing("source", "destination", generatorMock, txMock, nil)

	// When: calling ping
	timing, err := replPing.Ping(context.TODO(), "test")

	// Then: generated packet should be sent via txChannel and received via rxChannel
	require.NoError(t, err)
	// And: send & receive phases should be reported in timing record
	require.True(t, timing.GetPhase("send") > 0)
	require.True(t, timing.GetPhase("receive") > 0)
}

func TestReplicationPongShouldSendGeneratedPacketViaRxAndReceiveItViaRxOnPong(t *testing.T) {
	// Given: replication ping built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	packet := model.NewPingPacket("test", "destination", time.Second, 1)
	generatorMock := mocks.NewMockGenerator(ctrl)
	generatorMock.EXPECT().Generate("test", "destination").Return(packet)

	rxMock := mocks.NewMockChannel(ctrl)
	rxMock.EXPECT().Receive("test").Return(packet, nil)
	rxMock.EXPECT().Send(packet).Return(nil)
	replPing := pinger.NewReplicationPing("source", "destination", generatorMock, nil, rxMock)

	// When: calling pong
	timing, err := replPing.Pong(context.TODO(), "test")

	// Then: generated packet should be sent via txChannel and received via rxChannel
	require.NoError(t, err)
	// And: send & receive phases should be reported in timing record
	require.True(t, timing.GetPhase("send") > 0)
	require.True(t, timing.GetPhase("receive") > 0)
}

func TestReplicationHalfPingShouldSendGeneratedPacketOnPing(t *testing.T) {
	// Given: replication ping built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	packet := model.NewPingPacket("test", "source", time.Second, 1)
	generatorMock := mocks.NewMockGenerator(ctrl)
	generatorMock.EXPECT().Generate("test", "source").Return(packet)
	channelMock := mocks.NewMockChannel(ctrl)
	channelMock.EXPECT().Send(packet).Return(nil)

	replHalfPing := pinger.NewReplicationHalfPing("source", generatorMock, channelMock)

	// When: calling ping
	timing, err := replHalfPing.Ping(context.TODO(), "test")

	// Then: generated packet should be sent via channel
	require.NoError(t, err)
	// And: send phase should be reported in timing record
	require.True(t, timing.GetPhase("send") > 0)
	require.False(t, timing.GetPhase("receive") > 0)
}

func TestReplicationHalfPingShouldReceivePacketOnPong(t *testing.T) {
	// Given: replication ping built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	packet := model.NewPingPacket("test", "destination", time.Second, 1)
	channelMock := mocks.NewMockChannel(ctrl)
	channelMock.EXPECT().Receive("test").Return(packet, nil)

	replHalfPing := pinger.NewReplicationHalfPing("source", nil, channelMock)

	// When: calling pong
	timing, err := replHalfPing.Pong(context.TODO(), "test")

	// Then: packet should be received via channel
	require.NoError(t, err)
	// And: receive phase should be reported in timing record
	require.True(t, timing.GetPhase("receive") > 0)
	require.False(t, timing.GetPhase("send") > 0)
}
