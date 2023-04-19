package pinger_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces/mocks"
	"testing"
	"time"
)

func TestReplicationPingShouldSendGeneratedPacketViaTxAndReceiveItViaRxOnPing(t *testing.T) {
	// Given: replication ping built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	packet := pinger.NewPingPacket("test", "source", time.Second, 1)
	generatorMock := mocks.NewMockGenerator(ctrl)
	generatorMock.EXPECT().Generate("test", "source").Return(packet)
	txMock := mocks.NewMockChannel(ctrl)
	txMock.EXPECT().Send(packet).Return(nil)
	rxMock := mocks.NewMockChannel(ctrl)
	rxMock.EXPECT().Receive("test").Return(packet, nil)
	replPing := pinger.NewReplicationPing("source", "destination", generatorMock, txMock, rxMock)

	// When: calling ping
	err := replPing.Ping(context.TODO(), "test")

	// Then: generated packet should be send via txChannel and received via rxChannel
	require.NoError(t, err)
}

func TestReplicationPongShouldSendGeneratedPacketViaRxAndReceiveItViaTxOnPong(t *testing.T) {
	// Given: replication ping built on top of mocked dependencies
	ctrl := gomock.NewController(t)
	packet := pinger.NewPingPacket("test", "destination", time.Second, 1)
	generatorMock := mocks.NewMockGenerator(ctrl)
	generatorMock.EXPECT().Generate("test", "destination").Return(packet)

	txMock := mocks.NewMockChannel(ctrl)
	txMock.EXPECT().Receive("test").Return(packet, nil)
	rxMock := mocks.NewMockChannel(ctrl)
	rxMock.EXPECT().Send(packet).Return(nil)
	replPing := pinger.NewReplicationPing("source", "destination", generatorMock, txMock, rxMock)

	// When: calling pong
	err := replPing.Pong(context.TODO(), "test")

	// Then: generated packet should be send via txChannel and received via rxChannel
	require.NoError(t, err)
}
