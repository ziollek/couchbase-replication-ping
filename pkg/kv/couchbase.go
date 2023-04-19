package kv

import (
	"github.com/ziollek/couchbase-replication-ping/pkg/config"
	"github.com/ziollek/couchbase-replication-ping/pkg/kv/interfaces"
	"time"

	"github.com/couchbase/gocb/v2"
)

const GenericTimeout = 15 * time.Second

type Couchbase struct {
	collection *gocb.Collection
}

func buildConnectionOptions(c *config.Couchbase) gocb.ClusterOptions {
	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: c.User,
			Password: c.Password,
		},
	}
	// TODO: make timeouts configurable
	options.TimeoutsConfig.ConnectTimeout = GenericTimeout
	options.TimeoutsConfig.KVTimeout = GenericTimeout
	return options
}

func NewKV(c *config.Couchbase) (interfaces.KV, error) {
	cluster, err := gocb.Connect(c.URI, buildConnectionOptions(c))
	if err != nil {
		return nil, err
	}
	bucket := cluster.Bucket(c.Bucket)
	err = bucket.WaitUntilReady(GenericTimeout, nil)
	if err != nil {
		return nil, err
	}
	return &Couchbase{collection: bucket.DefaultCollection()}, nil
}

func (c *Couchbase) Upsert(key string, value interface{}, expiry time.Duration) error {
	_, err := c.collection.Upsert(key, value, &gocb.UpsertOptions{Expiry: expiry})
	return err
}

func (c *Couchbase) Get(key string, value interface{}) error {
	result, err := c.collection.Get(key, nil)
	if err != nil {
		return err
	}
	err = result.Content(value)
	return err
}
