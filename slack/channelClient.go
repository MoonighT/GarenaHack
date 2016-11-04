package slack

import (
	"fmt"
	"time"

	redis "gopkg.in/redis.v3"
)

type ChannelClient struct {
	client *redis.Client
}

func (c *ChannelClient) Init(addr string, db int) {
	options := &redis.Options{
		Network:      "tcp",
		Addr:         addr,
		DB:           int64(db),
		MaxRetries:   5,
		DialTimeout:  time.Duration(500) * time.Millisecond,
		ReadTimeout:  time.Duration(200) * time.Millisecond,
		WriteTimeout: time.Duration(500) * time.Millisecond,
		PoolSize:     10,
		PoolTimeout:  time.Duration(1) * time.Second,
		IdleTimeout:  time.Duration(1) * time.Minute,
	}
	c.client = redis.NewClient(options)
}

const (
	CHANNEL_KEY = "channelseq"
)

func (c *ChannelClient) Increase(channelid string) (seqid int64, err error) {
	field := fmt.Sprintf("%s", channelid)
	cmd := c.client.HIncrBy(CHANNEL_KEY, field, int64(1))
	if cmd.Err() != nil {
		return int64(-1), cmd.Err()
	}
	return cmd.Val(), cmd.Err()
}
