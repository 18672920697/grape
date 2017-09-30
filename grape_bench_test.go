package main

import (
	"github.com/go-redis/redis"
	"strconv"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:9001",
		Password: "",
		DB:       0,
	})

	for i := 0; i < 100; i++ {
		err := client.Set("key"+strconv.Itoa(i), "value", 0).Err()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:9001",
		Password: "",
		DB:       0,
	})

	for i := 0; i < 100; i++ {
		_, err := client.Get("key" + strconv.Itoa(i)).Result()
		if err != nil {
			panic(err)
		}
	}
}
