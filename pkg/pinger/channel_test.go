package pinger_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	kv_mocks "github.com/ziollek/couchbase-replication-ping/pkg/kv/interfaces/mocks"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger"
	"github.com/ziollek/couchbase-replication-ping/pkg/pinger/interfaces/mocks"
	"testing"
	"time"
)

func TestChannelShouldStorePacketInSourceStorageUnderProperKeyAndWithProperTTL(t *testing.T) {
	key := "test"
	ttl := time.Second * 5
	// Given: mocked packet & channel build on top of mocked dependencies
	ctrl := gomock.NewController(t)
	packet := mocks.NewMockPacket(ctrl)
	packet.EXPECT().GetKey().Return(key)
	packet.EXPECT().GetTTL().Return(ttl)

	kvMock := kv_mocks.NewMockKV(ctrl)
	kvMock.EXPECT().Upsert(key, packet, ttl)

	channel := pinger.NewChannel(kvMock, nil)

	// When: calling send
	err := channel.Send(packet)

	// Then: packet should be stored on source storage
	require.NoError(t, err)
}

func TestChannelShouldReadPacketFromDestinationStorageOnReceive(t *testing.T) {
	key := "test"
	// Given: channel build on top of mocked dependencies
	ctrl := gomock.NewController(t)

	kvMock := kv_mocks.NewMockKV(ctrl)
	kvMock.EXPECT().Get(key, gomock.Any()).Return(nil)

	channel := pinger.NewChannel(nil, kvMock)

	// When: calling send
	_, err := channel.Receive(key)

	// Then: packet should be read from destination storage
	require.NoError(t, err)
}
