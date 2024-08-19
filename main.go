package main

import (
	"SDC/consumer"
	"SDC/producer"
	"SDC/queue"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	PUB_SUB_QUEUE   = "pub_subs_queue"
	PUB_SUB_CHANNEL = "pub_sub_channel"
	NUM_CONSUMER    = 3
)

func bbq_demonstration() {
	queue := queue.NewBlockingBoundedQueue(5)

	go func() {
		for i := 0; i < 10; i++ {
			queue.Enqueue(i)
			fmt.Printf("Produced: %d\n", i)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			item := queue.Dequeue()
			fmt.Printf("Consumed: %v\n", item)
		}
	}()
}

func pub_sub_demo() {
	ctx, cancel := context.WithCancel(context.Background())
	rdb := setup_redis()

	defer cancel()

	var wg sync.WaitGroup

	numConsumers := NUM_CONSUMER

	for i := 1; i <= numConsumers; i++ {
		wg.Add(1)
		rc := consumer.NewRedisConsumer(rdb, PUB_SUB_QUEUE, PUB_SUB_CHANNEL, int64(i))
		go rc.Listener(ctx, &wg)
	}

	fmt.Printf("Created %v subscribers\n", NUM_CONSUMER)

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer_demo(rdb, ctx)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	fmt.Println("Shutting down gracefully...")
	cancel()
	wg.Wait()
}

func setup_redis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func producer_demo(rdb *redis.Client, ctx context.Context) {
	rp := producer.NewRedisProducer(rdb, PUB_SUB_QUEUE, PUB_SUB_CHANNEL)

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Message %d", i)
		rp.Publish(ctx, msg)
		if i%3 == 0 {
			time.Sleep(5 * time.Second)
		}
		time.Sleep(1 * time.Second)
	}
}

// func consumer_demo(rdb *redis.Client, ctx context.Context) {
// 	rc := consumer.NewRedisConsumer(rdb, PUB_SUB_QUEUE)

// 	rc.PushListener(ctx)
// }

func main() {

	pub_sub_demo()

	select {}
}
