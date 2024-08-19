package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConsumer struct {
	client  *redis.Client
	queue   string
	channel string
	id      int64
}

func NewRedisConsumer(rclient *redis.Client, queue string, channel string, id int64) *RedisConsumer {
	return &RedisConsumer{
		client:  rclient,
		queue:   queue,
		channel: channel,
		id:      id,
	}
}

func (rc *RedisConsumer) Listener(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer received cancellation signal, stopping listener...")
			return
		default:
			msg, err := rc.client.RPop(ctx, rc.queue).Result()
			if err == redis.Nil {
				fmt.Println("Queue is empty, waiting...")
				time.Sleep(1 * time.Second)
			} else if err != nil {
				log.Fatalf("Could not pop from queue: %v", err)
			} else {
				fmt.Printf("[Consumer%v] :: Consumed: %s\n", rc.id, msg)
			}
		}
	}
}

func (rc *RedisConsumer) PushListener(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	pubsub := rc.client.Subscribe(ctx, rc.channel)

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Could not subscribe: %v\n", err)
	}

	ch := pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			fmt.Printf("Subscriber %d: Received %s\n", rc.id, msg.Payload)
		case <-ctx.Done():
			fmt.Printf("Subscriber %d: Received cancellation signal, stopping listener...\n", rc.id)
			return
		}
	}
}
