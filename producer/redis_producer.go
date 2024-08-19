package producer

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisProducer struct {
	client  *redis.Client
	queue   string
	channel string
}

func NewRedisProducer(client *redis.Client, queue string, channel string) *RedisProducer {
	return &RedisProducer{client: client, queue: queue, channel: channel}
}

func (rp *RedisProducer) Publish(ctx context.Context, item interface{}) {
	err := rp.client.LPush(ctx, rp.queue, item).Err()
	if err != nil {
		log.Fatalf("Could not push to queue: %v", err)
	}
	log.Default().Printf("Produced: %s\n", item)
}

func (rp *RedisProducer) PushPublish(ctx context.Context, item interface{}) {
	err := rp.client.Publish(ctx, rp.channel, item).Err()
	if err != nil {
		log.Fatalf("Could not push to queue: %v", err)
	}
	log.Default().Printf("Produced: %s\n", item)
}
