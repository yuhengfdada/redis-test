package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	c := NewFailoverClient()
	err := c.FlushAll(c.Context()).Err()
	if err != nil {
		panic(err)
	}
	err = c.Set(c.Context(), "k1", "v1", 0).Err()
	if err != nil {
		panic(err)
	}
	err = c.Set(c.Context(), "k2", "v2", 0).Err()
	if err != nil {
		panic(err)
	}

	roc := NewFailoverReadOnlyClient()

	for i := 0; i < 3; i++ {
		res, err := roc.Get(roc.Context(), "k1").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
	fmt.Println("going to sleep")
	time.Sleep(3 * time.Second)
	fmt.Println("wake up")
	currentTime := time.Now()
	fmt.Println("Wake up time: ", currentTime.Format("2006.01.02 15:04:05"))

	for i := 0; i < 3; i++ {
		res, err := roc.Get(roc.Context(), "k2").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}

}

func NewFailoverClient() *redis.Client {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "myprimary",
		SentinelAddrs: []string{"127.0.0.1:5000", "localhost:5001", "localhost:5002"},
		DB:            0,
		Password:      "password",
	})
}

func NewFailoverReadOnlyClient() *redis.Client {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "myprimary",
		SentinelAddrs: []string{"localhost:5000", "localhost:5001", "localhost:5002"},
		SlaveOnly:     true,
		DB:            0,
		Password:      "password",
	})
}
